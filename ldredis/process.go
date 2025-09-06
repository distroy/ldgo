/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/distroy/ldgo/v3/ldredis/internal"
	redis "github.com/redis/go-redis/v9"
)

func getCallerField(rds *Redis) ldlog.Attr {
	return internal.GetCallerField(rds.opts.Caller)
}

func getCmdField(cmd Cmder) ldlog.Attr {
	return ldlog.Reflect("cmd", cmd.Args())
}

func newHook(rds *Redis) redis.Hook {
	return hook{
		Redis: rds,
	}
}

type hook struct {
	Redis *Redis
}

func (h hook) DialHook(next redis.DialHook) redis.DialHook { return nil }
func (h hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(c context.Context, cmd redis.Cmder) error {
		return h.process(c, cmd, next)
	}
}
func (h hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(c context.Context, cmds []redis.Cmder) error {
		return h.processPipeline(c, cmds, next)
	}
}

func (h hook) process(c context.Context, cmd Cmder, next redis.ProcessHook) error {
	if internal.InProcess(c) {
		return next(c, cmd)
	}

	var (
		ctx           = internal.NewContext(c)
		retry         = h.Redis.opts.Retry
		retryInterval = h.Redis.opts.RetryInterval
		reporter      = h.Redis.opts.Reporter
		logger        = ldctx.GetLogger(ctx)
		caller        = getCallerField(h.Redis)
	)

	for i := 0; ; {
		begin := time.Now()
		err := next(ctx, cmd)
		cmd.SetErr(err)
		err = h.getCmdError(c, cmd)
		reporter.Report(cmd, time.Since(begin))
		if isErrNil(err) {
			logger.Debug("redis cmd succ", ldlog.Int("retry", i), getCmdField(cmd), caller)
			return err
		}

		i++
		if i >= retry {
			logger.Error("redis cmd fail", ldlog.Int("retry", i), getCmdField(cmd),
				ldlog.Error(err), caller)
			return err
		}

		time.Sleep(retryInterval)
	}
}

func (h hook) processPipeline(c context.Context, cmds []Cmder, next redis.ProcessPipelineHook) error {
	if internal.InProcess(c) {
		return next(c, cmds)
	}
	var (
		ctx           = internal.NewContext(c)
		retry         = h.Redis.opts.Retry
		retryInterval = h.Redis.opts.RetryInterval
		reporter      = h.Redis.opts.Reporter
		logger        = ldctx.GetLogger(ctx)
		caller        = getCallerField(h.Redis)
	)
	logger = logger.With(ldlog.String("pipeline", hex.EncodeToString(ldrand.Bytes(8))))

	for i := 0; ; {
		begin := time.Now()
		next(ctx, cmds) // nolint
		reporter.ReportPipeline(cmds, time.Since(begin))

		err := h.checkPipelineError(ctx, cmds)
		if isErrNil(err) {
			h.printPipelineSuccLog(cmds, i, logger, caller)
			logger.Debug("redis pipeline cmd succ", ldlog.Int("retry", i), caller)
			return err
		}

		i++
		if i >= retry {
			h.printPipelineFailLog(cmds, i, logger, caller)
			logger.Error("redis pipeline fail", ldlog.Int("retry", i), ldlog.Error(err), caller)
			return err
		}

		time.Sleep(retryInterval)
	}
}

func (h hook) getCmdError(c context.Context, cmd Cmder) error {
	err := cmd.Err()
	if !isErrNil(err) {
		return err
	}

	v, _ := cmd.(internal.CmderWithParse)
	if v == nil {
		return err
	}

	err = v.Parse(c)
	if !isErrNil(err) {
		v.SetErr(err)
		return err
	}
	return err
}

func (h hook) checkPipelineError(c context.Context, cmds []Cmder) error {
	for _, cmd := range cmds {
		err := h.getCmdError(c, cmd)
		if !isErrNil(err) {
			return err
		}
	}
	return nil
}

func (h hook) printPipelineSuccLog(cmds []Cmder, retry int, logger *ldlog.Logger, caller ldlog.Attr) {
	for _, cmd := range cmds {
		logger.Debug("redis pipeline cmd succ", ldlog.Int("retry", retry), getCmdField(cmd), caller)
	}
}

func (h hook) printPipelineFailLog(cmds []Cmder, retry int, logger *ldlog.Logger, caller ldlog.Attr) {
	for _, cmd := range cmds {
		if err := cmd.Err(); !isErrNil(err) {
			logger.Error("redis pipeline cmd fail", ldlog.Int("retry", retry), getCmdField(cmd),
				ldlog.Error(err), caller)
			break
		}
		logger.Debug("redis pipeline cmd succ", ldlog.Int("retry", retry), getCmdField(cmd), caller)
	}
}

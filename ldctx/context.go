/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"context"
	"reflect"
	"time"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/distroy/ldgo/v2/ldlog"
)

type (
	Context    = context.Context
	CancelFunc = context.CancelFunc
)

func Default() context.Context { return defaultContext }
func Console() context.Context { return consoleContext }
func Discard() context.Context { return discardContext }

func ContextName(c context.Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflect.TypeOf(c).String()
}

func GetError(c context.Context) error {
	e := c.Err()
	switch e {
	case nil:
		return nil

	case context.Canceled:
		return lderr.ErrCtxCanceled

	case context.DeadlineExceeded:
		return lderr.ErrCtxDeadlineExceeded
	}

	return e
}

func WithLogger(c context.Context, log *ldlog.Logger, fields ...ldlog.Field) context.Context {
	if log == nil {
		return WithLogField(c, fields...)
	}
	return ctxWithLogger(c, func(_ *ldlog.Logger) *ldlog.Logger {
		return log.With(fields...)
	})
}

func WithLogField(c context.Context, fields ...ldlog.Field) context.Context {
	if len(fields) == 0 {
		return c
	}
	return ctxWithLogger(c, func(log *ldlog.Logger) *ldlog.Logger {
		return log.With(fields...)
	})
}

func WithLogEnabler(c context.Context, enabler ldlog.Enabler) context.Context {
	return ctxWithLogger(c, func(log *ldlog.Logger) *ldlog.Logger {
		return log.WithEnabler(enabler)
	})
}

// Log based on probability(rate). rate should be in [0, 1.0]
//
// Deprecated: use `WithLogEnabler` instead.
func WithLogRate(c context.Context, rate float64) context.Context {
	return ctxWithLogger(c, func(log *ldlog.Logger) *ldlog.Logger {
		return log.WithRate(rate)
	})
}

// Log based on time interval.
//
// Deprecated: use `WithLogEnabler` instead.
func WithLogInterval(c context.Context, interval time.Duration) context.Context {
	return ctxWithLogger(c, func(log *ldlog.Logger) *ldlog.Logger {
		return log.WithInterval(interval)
	})
}

func WithSequence(c context.Context, seq string) context.Context {
	if seq == "" {
		return c
	}
	return ctxWithLogger(c, func(log *ldlog.Logger) *ldlog.Logger {
		return log.WithSequence(seq)
	})
}

func GetSequence(c context.Context) string { return GetLogger(c).GetSequence() }

func ctxWithLogger(c context.Context, fnLog func(log *ldlog.Logger) *ldlog.Logger) context.Context {
	old := GetLogger(c)
	log := fnLog(old)
	if old == log {
		return c
	}
	return WithValue(c, ctxKeyLogger, log)
}

func GetLogger(c context.Context) *ldlog.Logger {
	if c == nil {
		return defaultLogger()
	}
	log, ok := c.Value(ctxKeyLogger).(*ldlog.Logger)
	if !ok {
		log = defaultLogger()
	}
	return log
}

func WithValue(parent context.Context, key, val interface{}) context.Context {
	return context.WithValue(parent, key, val)
}

func WithCancel(parent context.Context) (context.Context, CancelFunc) {
	return context.WithCancel(parent)
}

func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithDeadline(parent context.Context, deadline time.Time) (context.Context, CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

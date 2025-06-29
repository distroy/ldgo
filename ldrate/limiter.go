/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldatomic"
	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/distroy/ldgo/v3/ldmath"
	"golang.org/x/time/rate"
)

type Limiter struct {
	limiter *rate.Limiter
	config  config

	lastBurst     ldatomic.Int64
	lastLimit     ldatomic.Int64
	lastInterval  ldatomic.Duration
	lastNodeCount ldatomic.Int64
}

func NewLimiter(opts ...Option) *Limiter {
	cfg := newConfig(opts...)
	l := &Limiter{
		limiter: &rate.Limiter{},
		config:  cfg,
	}

	return l
}

func (l *Limiter) SetOptions(opts ...Option) {
	cfg := &l.config
	for _, fn := range opts {
		fn(cfg)
	}
}

func (l *Limiter) SetName(name string)         { l.config.Name.Store(name) }
func (l *Limiter) SetBurst(burst int64)        { l.config.Burst.Store(burst) }
func (l *Limiter) SetLimit(limit int64)        { l.config.Limit.Store(limit) }
func (l *Limiter) SetInterval(d time.Duration) { l.config.Interval.Store(d) }
func (l *Limiter) SetNodeCount(n int64)        { l.config.NodeCount.Store(n) }

func (l *Limiter) Name() string            { return l.config.Name.Load() }
func (l *Limiter) Burst() int64            { return l.config.Burst.Load() }
func (l *Limiter) Limit() int64            { return l.config.Limit.Load() }
func (l *Limiter) Interval() time.Duration { return l.config.Interval.Load() }
func (l *Limiter) NodeCount() int64        { return l.config.NodeCount.Load() }

func (l *Limiter) refresh(c context.Context, now time.Time) error {
	cfg := &l.config

	burst := cfg.Burst.Load()
	if burst <= 0 {
		burst = ldmath.Max(l.lastBurst.Load(), 1)
	}
	if burst != l.lastBurst.Load() {
		l.lastBurst.Store(burst)
		l.limiter.SetBurstAt(now, int(burst))
		ldctx.LogI(c, "[limiter] refresh the burst succ", ldlog.String("name", l.Name()),
			ldlog.Int64("burst", burst))
	}

	limit := cfg.Limit.Load()
	interval := cfg.Interval.Load()
	nodeCount := cfg.NodeCount.Load()

	if limit <= 0 {
		limit = ldmath.Max(l.lastLimit.Load(), 1)
	}
	if interval < 0 {
		interval = 0
	}
	if nodeCount <= 0 {
		nodeCount = ldmath.Max(l.lastNodeCount.Load(), 1)
	}

	// if interval < 0 || limit <= 0 || nodeCount <= 0 {
	// 	ldctx.LogE(c, "[limiter] invalid rate every parameters", ldlog.Int64("limit", limit),
	// 		ldlog.Stringer("interval", interval), ldlog.Int64("nodeCount", nodeCount))
	// 	return lderr.ErrInternalServerError
	// }

	if limit == l.lastLimit.Load() && interval == l.lastInterval.Load() && nodeCount == l.lastNodeCount.Load() {
		return nil
	}

	every := interval * time.Duration(nodeCount) / time.Duration(limit)
	// if every < 0 {
	// 	ldctx.LogE(c, "[limiter] invalid rate every", ldlog.Int64("limit", limit),
	// 		ldlog.Stringer("interval", interval), ldlog.Int64("serviceCount", nodeCount))
	// 	return lderr.ErrInternalServerError
	// }

	l.limiter.SetLimitAt(now, rate.Every(every))
	l.lastLimit.Store(limit)
	l.lastInterval.Store(interval)
	l.lastNodeCount.Store(nodeCount)

	ldctx.LogI(c, "[limiter] refresh rate every succ", ldlog.String("name", l.Name()),
		ldlog.Int64("limit", limit), ldlog.Stringer("interval", interval),
		ldlog.Int64("serviceCount", nodeCount), ldlog.Stringer("every", every))
	return nil
}

func (l *Limiter) Wait(c context.Context) error {
	return l.WaitN(c, 1)
}

func (l *Limiter) WaitN(c context.Context, n int) error {
	if err := wait(c, l, n); err != nil {
		ldctx.LogE(c, "[limiter] wait fail", ldlog.String("name", l.Name()), ldlog.Int("n", n), ldlog.Error(err))
		return err
	}
	return nil
}

func (l *Limiter) Allow(c context.Context) bool {
	return l.AllowN(c, 1)
}

func (l *Limiter) AllowN(c context.Context, n int) bool {
	return allow(c, l, n)
}

func (l *Limiter) Reserve(c context.Context) (*Reservation, error) {
	return l.ReserveN(c, 1)
}

func (l *Limiter) ReserveN(c context.Context, n int) (*Reservation, error) {
	now := time.Now()
	if err := l.refresh(c, now); err != nil {
		return nil, err
	}

	reservation := l.limiter.ReserveN(now, n)
	if !reservation.OK() {
		return nil, lderr.ErrCtxDeadlineNotEnough
	}

	r := &Reservation{
		reservations: []*rate.Reservation{reservation},
	}
	return r, nil
}

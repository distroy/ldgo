/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"golang.org/x/time/rate"
)

var (
	_ ILimiter = (*Limiter)(nil)
	_ ILimiter = (*Limiters)(nil)
)

type ILimiter interface {
	Wait(c context.Context) error
	WaitN(c context.Context, n int) error

	Allow(c context.Context) bool
	AllowN(c context.Context, n int) bool

	Reserve(c context.Context) (*Reservation, error)
	ReserveN(c context.Context, n int) (*Reservation, error)
}

type reserver interface {
	ReserveN(c context.Context, n int) (*Reservation, error)
}

func wait(c context.Context, l reserver, n int) error {
	select {
	case <-c.Done():
		// return lderr.ErrCtxCanceled
		return ldctx.GetError(c)
	default:
	}

	now := time.Now()

	// Determine wait limit
	waitLimit := rate.InfDuration
	if deadline, ok := c.Deadline(); ok {
		waitLimit = deadline.Sub(now)
	}

	r, err := l.ReserveN(c, n)
	if err != nil {
		return err
	}

	// Wait if necessary
	delay := r.DelayFrom(now)
	if delay <= 0 {
		return nil
	}

	if delay >= waitLimit {
		r.CancelAt(now)
		return lderr.ErrCtxDeadlineNotEnough
	}

	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		break

	case <-c.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.CancelAt(now)
		return ldctx.GetError(c)
	}

	// We can proceed.
	return nil
}

func allow(c context.Context, l reserver, n int) bool {
	r, err := l.ReserveN(c, n)
	if err != nil {
		return false
	}

	// Wait if necessary
	now := time.Now()
	delay := r.DelayFrom(now)
	// ldctx.LogI(c, " ** allow delay", zap.Duration("delay", delay))
	if delay <= 0 {
		return true
	}

	r.CancelAt(now)
	return false
}

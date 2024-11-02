/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"sync"
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/smartystreets/goconvey/convey"
)

type gopool struct {
	wg sync.WaitGroup
}

func (p *gopool) Go(fn func()) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		fn()
	}()
}

func (p *gopool) Wait() {
	p.wg.Wait()
}

func TestLimiter_Wait(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx := ldctx.Default()
		goes := &gopool{}

		begin := time.Now()
		interval := time.Second

		l := NewLimiter()
		l.SetName("wait")
		l.SetBurst(1)
		l.SetLimit(1)
		l.SetInterval(interval)

		c.Convey("wait times", func() {
			sleep := interval / 10

			for i := 0; i < 3; i++ {
				n := time.Duration(i+1) * interval
				goes.Go(func() {
					err := l.Wait(ctx)

					c.So(err, convey.ShouldBeNil)
					c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(n))
				})
				time.Sleep(sleep)
			}

			goes.Wait()
		})

		c.Convey("context has cancelled", func() {
			ctx, cancel := ldctx.WithCancel(ctx)
			cancel()

			err := l.Wait(ctx)

			c.So(err, convey.ShouldResemble, lderr.ErrCtxCanceled)
		})

		c.Convey("deedline not enough", func() {
			ctx, _ := ldctx.WithTimeout(ctx, interval/2)

			err := l.Wait(ctx)

			c.So(err, convey.ShouldResemble, lderr.ErrCtxDeadlineNotEnough)
		})

		c.Convey("no wait time", func() {
			l.refresh(ctx, begin.Add(-interval))
			// time.Sleep(interval)

			err := l.Wait(ctx)
			end := time.Now()
			c.So(err, convey.ShouldBeNil)
			c.So(end, convey.ShouldHappenBefore, begin.Add(interval))
			c.So(end, convey.ShouldHappenBefore, begin.Add(1*time.Millisecond))
		})
	})
}

func TestLimiter_Allow(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx := ldctx.Default()
		ctx, cancel := ldctx.WithCancel(ctx)
		cancel()

		interval := time.Second

		l := NewLimiter()
		l.SetName("allow")
		l.SetBurst(1)
		l.SetLimit(1)
		l.SetInterval(interval)
		c.So(l.Allow(ctx), convey.ShouldBeFalse)

		time.Sleep(interval * 2)
		c.So(l.Allow(ctx), convey.ShouldBeTrue)
		c.So(l.Allow(ctx), convey.ShouldBeFalse)

		time.Sleep(interval)
		c.So(l.Allow(ctx), convey.ShouldBeTrue)
		c.So(l.Allow(ctx), convey.ShouldBeFalse)

		time.Sleep(interval)
		c.So(l.Allow(ctx), convey.ShouldBeTrue)
		c.So(l.Allow(ctx), convey.ShouldBeFalse)
	})
}

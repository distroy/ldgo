/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"sync"
	"testing"
	"time"

	"github.com/distroy/ldgo/v2/ldctx"
	"github.com/distroy/ldgo/v2/lderr"
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
		l.SetBurst(1)
		l.SetLimit(1)
		l.SetInterval(interval)

		c.Convey("wait 10 times", func() {
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
			ctx := ldctx.WithCancel(ctx)
			ldctx.TryCancel(ctx)

			err := l.Wait(ctx)

			c.So(err, convey.ShouldEqual, lderr.ErrCtxCanceled)
		})

		c.Convey("deedline not enough", func() {
			ctx := ldctx.WithTimeout(ctx, interval/2)

			err := l.Wait(ctx)

			c.So(err, convey.ShouldEqual, lderr.ErrCtxDeadlineNotEnough)
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

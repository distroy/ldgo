/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/smartystreets/goconvey/convey"
)

func TestLimiters_Wait(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		goes := &gopool{}

		var (
			interval0 = time.Millisecond * 500
			interval1 = time.Millisecond * 700
			interval2 = time.Millisecond * 1000
		)

		var (
			l0 = NewLimiter(WithInterval(interval0))
			l1 = NewLimiter(WithInterval(interval1))
			l2 = NewLimiter(WithInterval(interval2))
		)

		var (
			ctx   = ldctx.Default()
			begin = time.Now()
			sleep = time.Millisecond * 10
		)

		l0.refresh(ctx, begin)
		l1.refresh(ctx, begin)
		l2.refresh(ctx, begin)

		var (
			l01 = NewLimiters(l0, l1)
			l02 = NewLimiters(l0, l2)
			l12 = NewLimiters(l1, l2)
		)

		goes.Go(func() {
			err := l01.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval1))
		})

		time.Sleep(sleep)
		goes.Go(func() {
			err := l12.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval1*2))
		})

		time.Sleep(sleep)
		goes.Go(func() {
			err := l02.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval2*2))
		})

		goes.Wait()
	})
}

func TestLimiters_Allow(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		var (
			interval0 = time.Millisecond * 500
			interval1 = time.Millisecond * 700
			interval2 = time.Millisecond * 1000
		)

		var (
			l0 = NewLimiter(WithInterval(interval0))
			l1 = NewLimiter(WithInterval(interval1))
			l2 = NewLimiter(WithInterval(interval2))
		)

		var (
			ctx   = ldctx.Default()
			begin = time.Now()
		)

		l0.refresh(ctx, begin)
		l1.refresh(ctx, begin)
		l2.refresh(ctx, begin)

		var (
			l01 = NewLimiters(l0, l1)
			l02 = NewLimiters(l0, l2)
			l12 = NewLimiters(l1, l2)
		)

		time.Sleep(700 * time.Millisecond)
		c.So(l01.Allow(ctx), convey.ShouldBeTrue)
		c.So(l01.Allow(ctx), convey.ShouldBeFalse)
		c.So(l02.Allow(ctx), convey.ShouldBeFalse)
		c.So(l12.Allow(ctx), convey.ShouldBeFalse)

		time.Sleep(500 * time.Millisecond)
		c.So(l02.Allow(ctx), convey.ShouldBeTrue)
		c.So(l01.Allow(ctx), convey.ShouldBeFalse)
		c.So(l02.Allow(ctx), convey.ShouldBeFalse)
		c.So(l12.Allow(ctx), convey.ShouldBeFalse)

		time.Sleep(1000 * time.Millisecond)
		c.So(l12.Allow(ctx), convey.ShouldBeTrue)
		c.So(l01.Allow(ctx), convey.ShouldBeFalse)
		c.So(l02.Allow(ctx), convey.ShouldBeFalse)
		c.So(l12.Allow(ctx), convey.ShouldBeFalse)
	})
}

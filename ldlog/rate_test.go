/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"io"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func Test_core_enable(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		test_core_enable_rate(c)
		test_core_enable_interval(c)
	})
}

func test_core_enable_rate(c convey.C) {
	l := newCore(newHandler(io.Discard))
	lvl := LevelInfo

	c.Convey("rate", func(c convey.C) {
		l.enabler = IntervalEnabler(time.Second)
		c.Convey("1", func(c convey.C) {
			l.enabler = RateEnabler(1)
			for range 100 {
				c.So(l.enabled(nil, lvl, 0), convey.ShouldBeTrue)
			}
		})
		c.Convey("0", func(c convey.C) {
			l.enabler = RateEnabler(0)
			for range 100 {
				c.So(l.enabled(nil, lvl, 0), convey.ShouldBeFalse)
			}
		})
		c.Convey("0.5", func(c convey.C) {
			l.enabler = RateEnabler(0.5)
			var (
				totalCnt = 20000
				diff     = 1000
			)
			trueCnt := 0
			for range totalCnt {
				if l.enabled(nil, lvl, 0) {
					trueCnt++
				}
			}
			half := totalCnt / 2
			c.So(trueCnt, convey.ShouldBeBetweenOrEqual, half-diff, half+diff)
		})
	})
}

func test_core_enable_interval(c convey.C) {
	l := newCore(newHandler(io.Discard))
	lvl := LevelInfo

	c.Convey("interval", func(c convey.C) {
		l.enabler = RateEnabler(0)
		c.Convey("0", func(c convey.C) {
			l.enabler = IntervalEnabler(0)
			for range 100 {
				c.So(l.enabled(nil, lvl, 0), convey.ShouldBeTrue)
			}
		})
		c.Convey("1s", func(c convey.C) {
			interval := time.Millisecond * 50
			l.enabler = IntervalEnabler(interval)

			time.Sleep(interval)
			c.So(l.enabled(nil, lvl, 1), convey.ShouldBeTrue)
			for range 100 {
				c.So(l.enabled(nil, lvl, 1), convey.ShouldBeFalse)
			}
			time.Sleep(interval)
			c.So(l.enabled(nil, lvl, 1), convey.ShouldBeTrue)
			c.So(l.enabled(nil, lvl, 1), convey.ShouldBeFalse)
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap/zapcore"
)

func TestRateEnabler(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("1", func(c convey.C) {
			l := RateEnabler(1)
			c.So(l, convey.ShouldResemble, defaultEnabler{})
		})
		c.Convey("0", func(c convey.C) {
			l := RateEnabler(0)
			c.So(l, convey.ShouldResemble, falseEnabler{})
		})
		c.Convey("0.5", func(c convey.C) {
			l := RateEnabler(0.5)
			c.So(l, convey.ShouldResemble, rateEnabler{rate: 0.5})

			testRateEnabler_HalfEnable(c, l)
		})
	})
}

func TestIntervalEnabler(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("0", func(c convey.C) {
			l := IntervalEnabler(0)
			c.So(l, convey.ShouldResemble, defaultEnabler{})
		})
		c.Convey("50ms", func(c convey.C) {
			const (
				info = zapcore.InfoLevel
				err  = zapcore.ErrorLevel

				interval = time.Millisecond * 50
			)

			l := IntervalEnabler(interval)
			c.So(l, convey.ShouldResemble, intervalEnabler{interval: interval})

			time.Sleep(interval)
			c.So(l.Enable(info, 1), convey.ShouldBeTrue)
			c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			for i := 0; i < 100; i++ {
				c.So(l.Enable(info, 1), convey.ShouldBeFalse)
				c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			}
			time.Sleep(interval)
			c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			c.So(l.Enable(info, 1), convey.ShouldBeTrue)
			c.So(l.Enable(info, 1), convey.ShouldBeFalse)
		})
	})
}

func Test_rateEnabler_Enable(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		l := &rateEnabler{}
		const (
			info = zapcore.InfoLevel
			err  = zapcore.ErrorLevel
		)

		c.Convey("1", func(c convey.C) {
			l.rate = 1
			for i := 0; i < 100; i++ {
				c.So(l.Enable(info), convey.ShouldBeTrue)
			}
		})
		c.Convey("0", func(c convey.C) {
			l.rate = 0
			for i := 0; i < 100; i++ {
				c.So(l.Enable(info), convey.ShouldBeFalse)
			}
			c.So(l.Enable(err), convey.ShouldBeTrue)
		})
		c.Convey("0.5", func(c convey.C) {
			l.rate = 0.5
			testRateEnabler_HalfEnable(c, l)
		})
	})
}

func testRateEnabler_HalfEnable(c convey.C, l Enabler) {
	const (
		info = zapcore.InfoLevel
		err  = zapcore.ErrorLevel
	)
	c.Convey("info log", func(c convey.C) {
		var (
			totalCnt = 20000
			diff     = 1000
		)
		trueCnt := 0
		for i := 0; i < totalCnt; i++ {
			if l.Enable(info) {
				trueCnt++
			}
		}
		half := totalCnt / 2
		c.So(trueCnt, convey.ShouldBeBetweenOrEqual, half-diff, half+diff)
	})
	c.Convey("error log", func(c convey.C) {
		for i := 0; i < 100; i++ {
			c.So(l.Enable(err), convey.ShouldBeTrue)
		}
	})
}

func Test_intervalEnabler_Enable(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		l := &intervalEnabler{}
		info := zapcore.InfoLevel
		err := zapcore.ErrorLevel

		c.Convey("0", func(c convey.C) {
			l.interval = 0
			for i := 0; i < 100; i++ {
				c.So(l.Enable(info, 0), convey.ShouldBeTrue)
				c.So(l.Enable(err, 0), convey.ShouldBeTrue)
			}
		})
		c.Convey("50ms", func(c convey.C) {
			interval := time.Millisecond * 50
			l.interval = interval

			time.Sleep(interval)
			c.So(l.Enable(info, 1), convey.ShouldBeTrue)
			c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			for i := 0; i < 100; i++ {
				c.So(l.Enable(info, 1), convey.ShouldBeFalse)
				c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			}
			time.Sleep(interval)
			c.So(l.Enable(err, 1), convey.ShouldBeTrue)
			c.So(l.Enable(info, 1), convey.ShouldBeTrue)
			c.So(l.Enable(info, 1), convey.ShouldBeFalse)
		})
	})
}

func TestEnablerByNameAndInterval(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		const (
			name0 = "test0"
			name1 = "test1"
		)
		c.Convey("0ms", func(c convey.C) {
			interval := time.Duration(0)

			e00 := EnablerByNameAndInterval(name0, interval)
			e01 := EnablerByNameAndInterval(name0, interval)
			e02 := EnablerByNameAndInterval(name0, interval)

			e10 := EnablerByNameAndInterval(name1, interval)
			e11 := EnablerByNameAndInterval(name1, interval)
			e12 := EnablerByNameAndInterval(name1, interval)

			c.So(e00, convey.ShouldResemble, defaultEnabler{})
			c.So(e01, convey.ShouldResemble, defaultEnabler{})
			c.So(e02, convey.ShouldResemble, defaultEnabler{})

			c.So(e10, convey.ShouldResemble, defaultEnabler{})
			c.So(e11, convey.ShouldResemble, defaultEnabler{})
			c.So(e12, convey.ShouldResemble, defaultEnabler{})
		})
		c.Convey("50ms", func(c convey.C) {
			interval := time.Millisecond * 50

			e00 := EnablerByNameAndInterval(name0, interval)
			e01 := EnablerByNameAndInterval(name0, interval)
			e02 := EnablerByNameAndInterval(name0, interval)

			e10 := EnablerByNameAndInterval(name1, interval)
			e11 := EnablerByNameAndInterval(name1, interval)
			e12 := EnablerByNameAndInterval(name1, interval)

			time.Sleep(interval)
			e03 := EnablerByNameAndInterval(name0, interval)
			e13 := EnablerByNameAndInterval(name1, interval)

			c.So(e00, convey.ShouldResemble, defaultEnabler{})
			c.So(e01, convey.ShouldResemble, falseEnabler{})
			c.So(e02, convey.ShouldResemble, falseEnabler{})
			c.So(e03, convey.ShouldResemble, defaultEnabler{})

			c.So(e10, convey.ShouldResemble, defaultEnabler{})
			c.So(e11, convey.ShouldResemble, falseEnabler{})
			c.So(e12, convey.ShouldResemble, falseEnabler{})
			c.So(e13, convey.ShouldResemble, defaultEnabler{})
		})
	})
}

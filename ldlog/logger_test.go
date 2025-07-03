/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"bytes"
	"testing"
	"time"
	"unsafe"

	"github.com/distroy/ldgo/v3/ldhook"
	"github.com/distroy/ldgo/v3/ldptr"
	"github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		patches := ldhook.NewPatches()
		defer patches.Reset()
		patches.Applys([]ldhook.Hook{
			ldhook.FuncHook{
				Target: time.Now,
				Double: ldhook.Values{time.Unix(1629610258, 0)},
			},
		})

		type LoggerValue struct {
			Name string
		}

		writer := bytes.NewBuffer(nil)
		l := New(writer)
		l = l.With(String("abc", "xxx"))

		c.Convey("error", func(c convey.C) {
			l.Error("error message")
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldlog/logger_test.go:38|error message,abc=xxx\n")
		})

		c.Convey("warn", func(c convey.C) {
			l.Warn("warn message", Int("int", 123))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldlog/logger_test.go:44|warn message,abc=xxx|int=123\n")
		})

		c.Convey("info", func(c convey.C) {
			l.Infoln("info message", (10 * time.Millisecond))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldlog/logger_test.go:50|info message 10ms,abc=xxx\n")
		})

		c.Convey("warnln", func(c convey.C) {
			l.Warnln("warnln message", ldptr.New(1234), map[string]string{"a": "b"})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldlog/logger_test.go:56|warnln message 1234 map[a:b],abc=xxx\n")
		})

		c.Convey("infoln", func(c convey.C) {
			l.Infoln("infoln message", &LoggerValue{Name: "abc"}, []any{ldptr.New(1234), (*int)(nil)})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldlog/logger_test.go:62|infoln message {Name:abc} [1234, (*int)(nil)],abc=xxx\n")
		})

		c.Convey("errorln", func(c convey.C) {
			l.Errorln("errorln message", (*LoggerValue)(nil), unsafe.Pointer(uintptr(0x2345)))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldlog/logger_test.go:68|errorln message (*ldlog.LoggerValue)(nil) (unsafe.Pointer)(0x2345),abc=xxx\n")
		})

		c.Convey("map", func(c convey.C) {
			l.Warnln("warnln message", ldptr.New(1234), map[any]any{
				"a":       "b",
				100:       124,
				int64(10): 234,
			})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldlog/logger_test.go:74|warnln message 1234 map[10:234,100:124,a:b],abc=xxx\n")
		})

		c.Convey("errorf", func(c convey.C) {
			l.Errorf("errorf message. int:%d", 1234)
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldlog/logger_test.go:84|errorf message. int:1234,abc=xxx\n")
		})
	})
}

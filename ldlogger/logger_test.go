/*
 * Copyright (C) distroy
 */

package ldlogger

import (
	"bytes"
	"testing"
	"time"
	"unsafe"

	"github.com/distroy/ldgo/ldhook"
	"github.com/distroy/ldgo/ldptr"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
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
		l := NewLogger(Writer(writer))
		l = With(l, zap.String("abc", "xxx"))

		convey.Convey("error", func() {
			l.Error("error message")
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldlogger/logger_test.go:39|error message,abc=xxx\n")
		})

		convey.Convey("warn", func() {
			l.Warn("warn message")
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldlogger/logger_test.go:45|warn message,abc=xxx\n")
		})

		convey.Convey("info", func() {
			l.Info("info message")
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldlogger/logger_test.go:51|info message,abc=xxx\n")
		})

		convey.Convey("warnln", func() {
			l.Warnln("warnln message", ldptr.NewInt(1234), map[string]string{"a": "b"})
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldlogger/logger_test.go:57|warnln message 1234 map[a:b],abc=xxx\n")
		})

		convey.Convey("infoln", func() {
			l.Infoln("infoln message", &LoggerValue{Name: "abc"}, []interface{}{ldptr.NewInt(1234), (*int)(nil)})
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldlogger/logger_test.go:63|infoln message {Name:abc} [1234, (*int)(nil)],abc=xxx\n")
		})

		convey.Convey("errorln", func() {
			l.Errorln("errorln message", (*LoggerValue)(nil), (unsafe.Pointer)((uintptr)(0x2345)))
			convey.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldlogger/logger_test.go:69|errorln message (*ldlogger.LoggerValue)(nil) (unsafe.Pointer)(0x2345),abc=xxx\n")
		})
	})
}

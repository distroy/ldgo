/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDuration(t *testing.T) {
	type Object struct {
		Timeout Duration `json:"timeout"`
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		d0 := Duration(123456789)
		d1 := d0.Duration()
		c.Convey("Duration", func(c convey.C) {
			c.So(d0, convey.ShouldEqual, d1)
		})
		c.Convey("Abs", func(c convey.C) {
			c.So(d0.Abs(), convey.ShouldEqual, d1.Abs())
		})
		c.Convey("Hours", func(c convey.C) {
			c.So(d0.Hours(), convey.ShouldEqual, d1.Hours())
		})
		c.Convey("Microseconds", func(c convey.C) {
			c.So(d0.Microseconds(), convey.ShouldEqual, d1.Microseconds())
		})
		c.Convey("Milliseconds", func(c convey.C) {
			c.So(d0.Milliseconds(), convey.ShouldEqual, d1.Milliseconds())
		})
		c.Convey("Minutes", func(c convey.C) {
			c.So(d0.Minutes(), convey.ShouldEqual, d1.Minutes())
		})
		c.Convey("Nanoseconds", func(c convey.C) {
			c.So(d0.Nanoseconds(), convey.ShouldEqual, d1.Nanoseconds())
		})
		c.Convey("Round", func(c convey.C) {
			c.So(d0.Round(Duration(time.Microsecond)), convey.ShouldEqual, d1.Round(time.Microsecond))
		})
		c.Convey("Seconds", func(c convey.C) {
			c.So(d0.Seconds(), convey.ShouldEqual, d1.Seconds())
		})
		c.Convey("String", func(c convey.C) {
			c.So(d0.String(), convey.ShouldEqual, d1.String())
		})
		c.Convey("Truncate", func(c convey.C) {
			c.So(d0.Truncate(Duration(time.Microsecond)), convey.ShouldEqual, d1.Truncate(time.Microsecond))
		})
		c.Convey("json marshal", func(c convey.C) {
			p := &Object{Timeout: 123456789}
			raw, err := json.Marshal(p)
			c.So(err, convey.ShouldBeNil)
			c.So(string(raw), convey.ShouldEqual, fmt.Sprintf(`{"timeout":"%s"}`, p.Timeout.get().String()))
		})
		c.Convey("json unmarshal", func(c convey.C) {
			c.Convey("string", func(c convey.C) {
				str := `{"timeout": "2h1m"}`
				p := &Object{}
				err := json.Unmarshal([]byte(str), p)
				c.So(err, convey.ShouldBeNil)
				c.So(p, convey.ShouldResemble, &Object{Timeout: Duration(time.Hour*2 + time.Minute*1)})
			})
			c.Convey("number", func(c convey.C) {
				str := `{"timeout": 12345}`
				p := &Object{}
				err := json.Unmarshal([]byte(str), p)
				c.So(err, convey.ShouldBeNil)
				c.So(p, convey.ShouldResemble, &Object{Timeout: 12345})
			})
		})
	})
}

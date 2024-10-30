/*
 * Copyright (C) distroy
 */

package ldmath

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAbs(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("int", func(c convey.C) {
			c.So(Abs(0), convey.ShouldEqual, 0)
			c.So(Abs(100), convey.ShouldEqual, 100)
			c.So(Abs(-100), convey.ShouldEqual, 100)
		})
		c.Convey("float32", func(c convey.C) {
			c.So(Abs[float32](0), convey.ShouldEqual, 0)
			c.So(Abs[float32](100), convey.ShouldEqual, 100)
			c.So(Abs[float32](-100), convey.ShouldEqual, 100)

			c.So(IsNaN(Abs[float32](0)), convey.ShouldResemble, false)
			c.So(IsNaN(Abs(NaN32())), convey.ShouldResemble, true)

			c.So(IsInf(Abs(Inf32(1)), 1), convey.ShouldEqual, true)
			c.So(IsInf(Abs(Inf32(-1)), 1), convey.ShouldEqual, true)
		})
		c.Convey("float64", func(c convey.C) {
			c.So(Abs[float64](0), convey.ShouldEqual, 0)
			c.So(Abs[float64](100), convey.ShouldEqual, 100)
			c.So(Abs[float64](-100), convey.ShouldEqual, 100)

			c.So(IsNaN(Abs[float64](0)), convey.ShouldResemble, false)
			c.So(IsNaN(Abs(NaN64())), convey.ShouldResemble, true)

			c.So(IsInf(Abs(Inf64(1)), -1), convey.ShouldEqual, false)
			c.So(IsInf(Abs(Inf64(-1)), -1), convey.ShouldEqual, false)
		})
	})
}

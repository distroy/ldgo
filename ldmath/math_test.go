/*
 * Copyright (C) distroy
 */

package ldmath

import (
	"testing"
	"time"

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

func TestMax(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Max(0, time.Second, time.Millisecond), convey.ShouldEqual, time.Second)

		c.So(Max[int](3, 4), convey.ShouldEqual, 4)
		c.So(Max[int8](3, 4), convey.ShouldEqual, int8(4))
		c.So(Max[int16](3, 4), convey.ShouldEqual, int16(4))
		c.So(Max[int32](3, 4), convey.ShouldEqual, int32(4))
		c.So(Max[int64](3, 4), convey.ShouldEqual, int64(4))

		c.So(Max[uint](3, 4), convey.ShouldEqual, 4)
		c.So(Max[uint8](3, 4), convey.ShouldEqual, uint8(4))
		c.So(Max[uint16](3, 4), convey.ShouldEqual, uint16(4))
		c.So(Max[uint32](3, 4), convey.ShouldEqual, uint32(4))
		c.So(Max[uint64](3, 4), convey.ShouldEqual, uint64(4))

		c.So(Max[float32](3, 4), convey.ShouldEqual, float32(4))
		c.So(Max[float64](3, 4), convey.ShouldEqual, float64(4))
	})
}

func TestMin(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Min(time.Second, time.Millisecond, 10), convey.ShouldEqual, 10)

		c.So(Min[int](4, 3), convey.ShouldEqual, 3)
		c.So(Min[int8](4, 3), convey.ShouldEqual, int8(3))
		c.So(Min[int16](4, 3), convey.ShouldEqual, int16(3))
		c.So(Min[int32](4, 3), convey.ShouldEqual, int32(3))
		c.So(Min[int64](4, 3), convey.ShouldEqual, int64(3))

		c.So(Min[uint](4, 3), convey.ShouldEqual, 3)
		c.So(Min[uint8](4, 3), convey.ShouldEqual, uint8(3))
		c.So(Min[uint16](4, 3), convey.ShouldEqual, uint16(3))
		c.So(Min[uint32](4, 3), convey.ShouldEqual, uint32(3))
		c.So(Min[uint64](4, 3), convey.ShouldEqual, uint64(3))

		c.So(Min[float32](4, 3), convey.ShouldEqual, float32(3))
		c.So(Min[float64](4, 3), convey.ShouldEqual, float64(3))
	})
}

func TestSum(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Sum(0), convey.ShouldEqual, 0)
		c.So(Sum(123, -115, 35), convey.ShouldEqual, 123-115+35)
	})
}

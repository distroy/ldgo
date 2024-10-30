/*
 * Copyright (C) distroy
 */

package ldmath

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMax(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Max(0, time.Second, time.Millisecond), convey.ShouldEqual, time.Second)

		convey.So(Max[int](3, 4), convey.ShouldEqual, 4)
		convey.So(Max[int8](3, 4), convey.ShouldEqual, int8(4))
		convey.So(Max[int16](3, 4), convey.ShouldEqual, int16(4))
		convey.So(Max[int32](3, 4), convey.ShouldEqual, int32(4))
		convey.So(Max[int64](3, 4), convey.ShouldEqual, int64(4))

		convey.So(Max[uint](3, 4), convey.ShouldEqual, 4)
		convey.So(Max[uint8](3, 4), convey.ShouldEqual, uint8(4))
		convey.So(Max[uint16](3, 4), convey.ShouldEqual, uint16(4))
		convey.So(Max[uint32](3, 4), convey.ShouldEqual, uint32(4))
		convey.So(Max[uint64](3, 4), convey.ShouldEqual, uint64(4))

		convey.So(Max[float32](3, 4), convey.ShouldEqual, float32(4))
		convey.So(Max[float64](3, 4), convey.ShouldEqual, float64(4))
	})
}

func TestMin(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Min(time.Second, time.Millisecond, 10), convey.ShouldEqual, 10)

		convey.So(Min[int](4, 3), convey.ShouldEqual, 3)
		convey.So(Min[int8](4, 3), convey.ShouldEqual, int8(3))
		convey.So(Min[int16](4, 3), convey.ShouldEqual, int16(3))
		convey.So(Min[int32](4, 3), convey.ShouldEqual, int32(3))
		convey.So(Min[int64](4, 3), convey.ShouldEqual, int64(3))

		convey.So(Min[uint](4, 3), convey.ShouldEqual, 3)
		convey.So(Min[uint8](4, 3), convey.ShouldEqual, uint8(3))
		convey.So(Min[uint16](4, 3), convey.ShouldEqual, uint16(3))
		convey.So(Min[uint32](4, 3), convey.ShouldEqual, uint32(3))
		convey.So(Min[uint64](4, 3), convey.ShouldEqual, uint64(3))

		convey.So(Min[float32](4, 3), convey.ShouldEqual, float32(3))
		convey.So(Min[float64](4, 3), convey.ShouldEqual, float64(3))
	})
}

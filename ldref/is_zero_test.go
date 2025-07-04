/*
 * Copyright (C) distroy
 */

package ldref

import (
	"log"
	"testing"
	"unsafe"

	"github.com/smartystreets/goconvey/convey"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func TestIsZero(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("invalid", func(c convey.C) {
			c.So(IsZero(nil), convey.ShouldEqual, true)
		})

		c.Convey("bool", func(c convey.C) {
			c.So(IsZero(false), convey.ShouldEqual, true)
			c.So(IsZero(true), convey.ShouldEqual, false)
		})

		c.Convey("int", func(c convey.C) {
			c.So(IsZero(int(0)), convey.ShouldEqual, true)
			c.So(IsZero(int(1)), convey.ShouldEqual, false)
		})

		c.Convey("uint", func(c convey.C) {
			c.So(IsZero(uint(0)), convey.ShouldEqual, true)
			c.So(IsZero(uint(1)), convey.ShouldEqual, false)
		})

		c.Convey("float64", func(c convey.C) {
			c.So(IsZero(float64(0)), convey.ShouldEqual, true)
			c.So(IsZero(float64(1)), convey.ShouldEqual, false)
		})

		c.Convey("*int", func(c convey.C) {
			n := 1
			c.So(IsZero((*int)(nil)), convey.ShouldEqual, true)
			c.So(IsZero((*int)(&n)), convey.ShouldEqual, false)
		})

		c.Convey("[]int", func(c convey.C) {
			c.So(IsZero(([]int)(nil)), convey.ShouldEqual, true)
			c.So(IsZero([]int{}), convey.ShouldEqual, false)
		})

		c.Convey("[...]int", func(c convey.C) {
			c.So(IsZero([...]int{}), convey.ShouldEqual, true)
			c.So(IsZero([...]int{0, 0, 1}), convey.ShouldEqual, false)
		})

		c.Convey("string", func(c convey.C) {
			c.So(IsZero(""), convey.ShouldEqual, true)
			c.So(IsZero("1"), convey.ShouldEqual, false)
		})

		c.Convey("complex128", func(c convey.C) {
			c.So(IsZero(complex(0, 0)), convey.ShouldEqual, true)
			c.So(IsZero(complex(0, 1)), convey.ShouldEqual, false)
			c.So(IsZero(complex(1, 0)), convey.ShouldEqual, false)
		})

		c.Convey("unsafe.Pointer", func(c convey.C) {
			n := 1
			c.So(IsZero((unsafe.Pointer)(nil)), convey.ShouldEqual, true)
			c.So(IsZero((unsafe.Pointer)(&n)), convey.ShouldEqual, false)
		})

		c.Convey("map[int]struct{}", func(c convey.C) {
			c.So(IsZero((map[int]struct{})(nil)), convey.ShouldEqual, true)
			c.So(IsZero(map[int]struct{}{}), convey.ShouldEqual, false)
		})

		c.Convey("func ()", func(c convey.C) {
			c.So(IsZero((func(c convey.C))(nil)), convey.ShouldEqual, true)
			c.So(IsZero(func(c convey.C) {}), convey.ShouldEqual, false)
		})

		c.Convey("chan struct{}", func(c convey.C) {
			c.So(IsZero((chan struct{})(nil)), convey.ShouldEqual, true)
			c.So(IsZero(make(chan struct{})), convey.ShouldEqual, false)
		})

		c.Convey("struct", func(c convey.C) {
			type testStruct = struct {
				Int int
			}
			c.So(IsZero(testStruct{}), convey.ShouldEqual, true)
			c.So(IsZero(testStruct{Int: 1}), convey.ShouldEqual, false)
		})
	})
}

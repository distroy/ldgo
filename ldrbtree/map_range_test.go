/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/smartystreets/goconvey/convey"
)

func TestMap_Range(t *testing.T) {
	const retry = 20
	fnDeleteAll := testMapDeleteAll

	convey.Convey(t.Name(), t, func() {
		m := testNewMap()

		convey.Convey("range", func() {
			it := m.Range()
			for _, d := range _nums {
				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(d, convey.ShouldEqual, it.Key())
				it.Next()
			}
			convey.So(it.HasNext(), convey.ShouldBeFalse)
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("search range", func() {
			for i := 0; i < retry; i++ {
				d := ldrand.Intn(_count)
				fnDeleteAll(m, d)

				m.Insert(d, 0)
				m.Insert(d, 1)
				m.Insert(d, 2)

				it := m.SearchRange(d)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 0)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 1)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 2)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeFalse)
			}
		})
	})
}

func TestMap_ReverseRange(t *testing.T) {
	const retry = 20
	fnDeleteAll := testMapRDeleteAll

	convey.Convey(t.Name(), t, func() {
		m := testNewMap()

		convey.Convey("range", func() {
			it := m.RRange()
			for i := len(_nums) - 1; i >= 0; i-- {
				n := _nums[i]
				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(n, convey.ShouldEqual, it.Key())
				it.Next()
			}
			convey.So(it.HasNext(), convey.ShouldBeFalse)
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("search range", func() {
			for i := 0; i < retry; i++ {
				d := ldrand.Intn(_count)
				fnDeleteAll(m, d)

				m.Insert(d, 0)
				m.Insert(d, 1)
				m.Insert(d, 2)

				it := m.RSearchRange(d)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 2)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 1)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Key(), convey.ShouldEqual, d)
				convey.So(it.Value(), convey.ShouldEqual, 0)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeFalse)
			}
		})
	})
}

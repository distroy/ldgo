/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/smartystreets/goconvey/convey"
)

func TestRBTree_Range(t *testing.T) {
	const retry = 20
	fnDeleteAll := testRBTreeDeleteAll

	convey.Convey(t.Name(), t, func() {
		rbtree := testNewRBTree()

		convey.Convey("range", func() {
			it := rbtree.Range()
			for _, n := range _nums {
				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, n)
				it.Next()
			}
			convey.So(it.HasNext(), convey.ShouldBeFalse)
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("search range", func() {
			for i := 0; i < retry; i++ {
				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.SearchRange(d)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeFalse)
			}
		})
	})
}

func TestRBTree_ReverseRange(t *testing.T) {
	const retry = 20
	fnDeleteAll := testRBTreeRDeleteAll

	convey.Convey(t.Name(), t, func() {
		rbtree := testNewRBTree()

		convey.Convey("range", func() {
			it := rbtree.RRange()
			for i := len(_nums) - 1; i >= 0; i-- {
				n := _nums[i]
				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, n)
				it.Next()
			}
			convey.So(it.HasNext(), convey.ShouldBeFalse)
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("search range", func() {
			for i := 0; i < retry; i++ {
				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.RSearchRange(d)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeTrue)
				convey.So(it.Get(), convey.ShouldEqual, d)
				convey.So(func() { it.Next() }, convey.ShouldNotPanic)

				convey.So(it.HasNext(), convey.ShouldBeFalse)
			}
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/smartystreets/goconvey/convey"
)

func TestMap_Insert(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		m := testNewMap()

		c.Convey("insert", func(c convey.C) {
			c.Convey("len", func(c convey.C) {
				c.So(m.Len(), convey.ShouldEqual, len(_nums))
			})

			c.Convey("check parent", func(c convey.C) {
				c.So(m.tree.root.checkParent(m.tree.sentinel), convey.ShouldBeTrue)
			})

			c.Convey("check color", func(c convey.C) {
				c.So(m.tree.root.checkColor(m.tree.sentinel), convey.ShouldBeTrue)
			})

			c.Convey("clear", func(c convey.C) {
				m.Clear()
				c.So(m.Len(), convey.ShouldEqual, 0)
				c.So(m.Begin(), convey.ShouldResemble, m.End())
			})

			c.Convey("count", func(c convey.C) {
				n := ldrand.Intn(_count)
				count := m.Count(n)
				m.Insert(n, 0)
				m.Insert(n, 1)
				m.Insert(n, 2)
				c.So(m.Count(n), convey.ShouldEqual, count+3)
			})
		})

		c.Convey("insert or search", func(c convey.C) {
			c.So(m.Len(), convey.ShouldEqual, len(_nums))

			c.Convey("insert when exists", func(c convey.C) {
				n := ldrand.Intn(_count)
				it := m.InsertOrSearch(n, -100)
				c.So(it.Key(), convey.ShouldEqual, n)
				c.So(it.Value(-100), convey.ShouldEqual, n)
				c.So(it.Value(), convey.ShouldEqual, -100)
				c.So(m.Len(), convey.ShouldEqual, len(_nums))
			})

			c.Convey("insert when not exists", func(c convey.C) {
				it := m.InsertOrSearch(-100, -100)
				c.So(it.Key(), convey.ShouldEqual, -100)
				c.So(it.Value(), convey.ShouldEqual, -100)
				c.So(m.Len(), convey.ShouldEqual, len(_nums)+1)
			})
		})

		c.Convey("insert or assign", func(c convey.C) {
			m := &Map[int, int]{}
			c.So(m.Len(), convey.ShouldEqual, 0)

			m.InsertOrAssign(100, 0)
			c.So(m.Len(), convey.ShouldEqual, 1)
			c.So(m.Begin().Key(), convey.ShouldResemble, 100)
			c.So(m.Begin().Value(), convey.ShouldResemble, 0)
			c.So(m.Begin().Next(), convey.ShouldResemble, m.End())

			m.InsertOrAssign(100, 2)
			c.So(m.Len(), convey.ShouldEqual, 1)
			c.So(m.Begin().Key(), convey.ShouldResemble, 100)
			c.So(m.Begin().Value(), convey.ShouldResemble, 2)
			c.So(m.Begin().Next(), convey.ShouldResemble, m.End())
		})
	})
}
func TestMap_DuplicateData(t *testing.T) {
	fnDeleteAll := testMapDeleteAll

	convey.Convey(t.Name(), t, func(c convey.C) {
		rbtree := testNewMap()

		c.Convey("search", func(c convey.C) {
			d := ldrand.Intn(_count)
			fnDeleteAll(rbtree, d)

			rbtree.Insert(d, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d, 2)

			it := rbtree.Search(d)

			c.So(it.Key(), convey.ShouldEqual, d)
			c.So(it.Next().Key(), convey.ShouldEqual, d)
			c.So(it.Next().Next().Key(), convey.ShouldEqual, d)

			c.So(it.Value(), convey.ShouldEqual, 0)
			c.So(it.Next().Value(), convey.ShouldEqual, 1)
			c.So(it.Next().Next().Value(), convey.ShouldEqual, 2)
		})

		c.Convey("lower bound", func(c convey.C) {
			d := ldrand.Intn(_count)
			fnDeleteAll(rbtree, d)

			rbtree.Insert(d, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d, 2)

			it := rbtree.LowerBound(d)

			c.So(it.Key(), convey.ShouldEqual, d)
			c.So(it.Next().Key(), convey.ShouldEqual, d)
			c.So(it.Next().Next().Key(), convey.ShouldEqual, d)

			c.So(it.Value(), convey.ShouldEqual, 0)
			c.So(it.Next().Value(), convey.ShouldEqual, 1)
			c.So(it.Next().Next().Value(), convey.ShouldEqual, 2)
		})

		c.Convey("upper bound", func(c convey.C) {
			d := ldrand.Intn(_count)

			rbtree.Insert(d-1, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d+1, 2)

			it := rbtree.UpperBound(d)

			c.So(it.Key(), convey.ShouldBeGreaterThan, d)
			c.So(it.Prev().Key(), convey.ShouldBeLessThanOrEqualTo, d)
			// c.So(it.Next().Next().Key(), convey.ShouldEqual, d)
		})
	})
}

func TestMap_ReverseDuplicateData(t *testing.T) {
	fnDeleteAll := testMapRDeleteAll

	convey.Convey(t.Name(), t, func(c convey.C) {
		rbtree := testNewMap()

		c.Convey("search", func(c convey.C) {
			d := ldrand.Intn(_count)
			fnDeleteAll(rbtree, d)

			rbtree.Insert(d, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d, 2)

			it := rbtree.RSearch(d)

			c.So(it.Key(), convey.ShouldEqual, d)
			c.So(it.Next().Key(), convey.ShouldEqual, d)
			c.So(it.Next().Next().Key(), convey.ShouldEqual, d)

			c.So(it.Value(), convey.ShouldEqual, 2)
			c.So(it.Next().Value(), convey.ShouldEqual, 1)
			c.So(it.Next().Next().Value(), convey.ShouldEqual, 0)
		})

		c.Convey("lower bound", func(c convey.C) {
			d := ldrand.Intn(_count)
			fnDeleteAll(rbtree, d)

			rbtree.Insert(d, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d, 2)

			it := rbtree.RLowerBound(d)

			c.So(it.Key(), convey.ShouldEqual, d)
			c.So(it.Next().Key(), convey.ShouldEqual, d)
			c.So(it.Next().Next().Key(), convey.ShouldEqual, d)

			c.So(it.Value(), convey.ShouldEqual, 2)
			c.So(it.Next().Value(), convey.ShouldEqual, 1)
			c.So(it.Next().Next().Value(), convey.ShouldEqual, 0)
		})

		c.Convey("upper bound", func(c convey.C) {
			d := ldrand.Intn(_count)

			rbtree.Insert(d-1, 0)
			rbtree.Insert(d, 1)
			rbtree.Insert(d+1, 2)

			it := rbtree.RUpperBound(d)

			c.So(it.Key(), convey.ShouldBeLessThan, d)
			c.So(it.Prev().Key(), convey.ShouldBeGreaterThanOrEqualTo, d)
			// c.So(it.Next().Next().Key(), convey.ShouldEqual, d)
		})
	})
}

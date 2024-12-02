/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sort"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMap(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("new map by nil", func(c convey.C) {
			m := NewMap[int, int](nil)

			c.So(m.Size(), convey.ShouldEqual, 0)
			c.So(m.Map(), convey.ShouldResemble, map[int]int{})

			v, loaded := m.Load(123)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)

			v, loaded = m.Swap(123, 100)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)

			c.So(m.Size(), convey.ShouldEqual, 1)
			c.So(m.Get(123), convey.ShouldEqual, 100)

			m.Set(123, 200)
			c.So(m.Size(), convey.ShouldEqual, 1)
			c.So(m.Get(123), convey.ShouldEqual, 200)

			m.Del(234)
			c.So(m.Size(), convey.ShouldEqual, 1)
		})

		c.Convey("new map by map", func(c convey.C) {
			m := NewMap(map[int]int{
				1000: 1234,
				2000: 1234,
			})

			c.So(m.Size(), convey.ShouldEqual, 2)
			c.So(m.Map(), convey.ShouldResemble, map[int]int{
				1000: 1234,
				2000: 1234,
			})

			v, loaded := m.Load(123)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)
			c.So(m.Has(123), convey.ShouldBeFalse)

			v, loaded = m.Swap(123, 100)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)
			c.So(m.Has(123), convey.ShouldBeTrue)

			c.So(m.Size(), convey.ShouldEqual, 3)
			c.So(m.Get(123), convey.ShouldEqual, 100)

			v, loaded = m.Swap(123, 200)
			c.So(loaded, convey.ShouldBeTrue)
			c.So(v, convey.ShouldEqual, 100)

			c.So(m.Size(), convey.ShouldEqual, 3)
			c.So(m.Get(123), convey.ShouldEqual, 200)

			m.Del(234)
			c.So(m.Size(), convey.ShouldEqual, 3)

			m.Set(234, 100)
			c.So(m.Size(), convey.ShouldEqual, 4)
			c.So(m.Get(234), convey.ShouldEqual, 100)

			m.Del(234)
			c.So(m.Size(), convey.ShouldEqual, 3)

			c.So(m.Map(), convey.ShouldResemble, map[int]int{
				1000: 1234,
				2000: 1234,
				123:  200,
			})

			keys := m.Keys()
			sort.Sort(sort.IntSlice(keys))
			c.So(keys, convey.ShouldResemble, []int{
				123, 1000, 2000,
			})

			{
				v, loaded := m.LoadOrStore(123, 100)
				c.So(loaded, convey.ShouldBeTrue)
				c.So(v, convey.ShouldEqual, 200)
				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 1234,
					2000: 1234,
					123:  200,
				})

				v, loaded = m.LoadOrStore(234, 100)
				c.So(loaded, convey.ShouldBeFalse)
				c.So(v, convey.ShouldEqual, 100)
				c.So(m.Size(), convey.ShouldEqual, 4)
				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 1234,
					2000: 1234,
					123:  200,
					234:  100,
				})
			}
		})
	})
}

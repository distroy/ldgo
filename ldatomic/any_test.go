/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAny(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("int", func(c convey.C) {
			v := &Any[int]{}
			c.So(v.Load(), convey.ShouldEqual, 0)

			c.So(v.CompareAndSwap(0, 1), convey.ShouldBeTrue)
			c.So(v.Load(), convey.ShouldEqual, 1)

			c.So(v.Swap(20), convey.ShouldEqual, 1)
			c.So(v.Load(), convey.ShouldEqual, 20)

			c.So(v.CompareAndSwap(1, 2), convey.ShouldBeFalse)
			c.So(v.Load(), convey.ShouldEqual, 20)
		})
		c.Convey("uint", func(c convey.C) {
			v := &Any[uint]{}
			c.So(v.Load(), convey.ShouldEqual, 0)

			v.Store(1)
			c.So(v.Load(), convey.ShouldEqual, 1)

			c.So(v.Swap(20), convey.ShouldEqual, 1)
			c.So(v.Load(), convey.ShouldEqual, 20)

			c.So(v.CompareAndSwap(0, 1), convey.ShouldBeFalse)
			c.So(v.Load(), convey.ShouldEqual, 20)

			c.So(v.CompareAndSwap(1, 2), convey.ShouldBeFalse)
			c.So(v.Load(), convey.ShouldEqual, 20)
		})
		c.Convey("float64", func(c convey.C) {
			v := &Any[float64]{}
			c.So(v.Load(), convey.ShouldEqual, 0)

			c.So(v.Swap(1), convey.ShouldEqual, 0)
			c.So(v.Load(), convey.ShouldEqual, 1)
		})
		c.Convey("map", func(c convey.C) {
			v := &Any[map[string]string]{}
			c.So(v.Load(), convey.ShouldBeNil)

			c.So(v.CompareAndSwap(nil, map[string]string{}), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, map[string]string{})

			c.So(v.CompareAndSwap(nil, map[string]string{"a": "1"}), convey.ShouldEqual, false)
			c.So(v.Load(), convey.ShouldResemble, map[string]string{})

			// c.So(v.CompareAndSwap(map[string]string{}, map[string]string{"a": "1"}), convey.ShouldEqual, false)
			// c.So(v.Load(), convey.ShouldResemble, map[string]string{})
			//
			// c.So(v.CompareAndSwap(v.Load(), map[string]string{"a": "1"}), convey.ShouldEqual, true)
			// c.So(v.Load(), convey.ShouldResemble, map[string]string{"a": "1"})
			c.So(v.CompareAndSwap(map[string]string{}, map[string]string{"a": "1"}), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, map[string]string{"a": "1"})

			c.So(v.CompareAndSwap(map[string]string{"a": "1"}, map[string]string{"b": "2"}), convey.ShouldEqual, false)
			c.So(v.Load(), convey.ShouldResemble, map[string]string{"a": "1"})

			c.So(v.CompareAndSwap(v.Load(), map[string]string{"b": "2"}), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, map[string]string{"b": "2"})
		})

		c.Convey("error", func(c convey.C) {
			v := &Any[error]{}
			c.So(v.Load(), convey.ShouldBeNil)

			var (
				error1 = fmt.Errorf("error1")
				error2 = fmt.Errorf("error2")
				error3 = fmt.Errorf("error3")
			)

			c.So(v.CompareAndSwap(nil, error1), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, error1)

			c.So(v.CompareAndSwap(nil, error2), convey.ShouldEqual, false)
			c.So(v.Load(), convey.ShouldResemble, error1)

			c.So(v.CompareAndSwap(error1, error3), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, error3)

			c.So(v.CompareAndSwap(fmt.Errorf("error3"), error2), convey.ShouldEqual, false)
			c.So(v.Load(), convey.ShouldResemble, error3)

			c.So(v.CompareAndSwap(v.Load(), error2), convey.ShouldEqual, true)
			c.So(v.Load(), convey.ShouldResemble, error2)
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldbuf

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewRing(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		b := NewRing[int](4)
		var (
			v  int
			ok bool
		)

		c.So(b.Cap(), convey.ShouldEqual, 4)
		c.So(b.Size(), convey.ShouldEqual, 0)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, false)
		c.So(v, convey.ShouldEqual, 0)

		b.Put(1)
		b.Put(2)
		b.Put(3)
		b.Put(4)
		b.Put(5)
		c.So(b.Size(), convey.ShouldEqual, 4)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, true)
		c.So(v, convey.ShouldEqual, 1)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, true)
		c.So(v, convey.ShouldEqual, 2)
		c.So(b.Size(), convey.ShouldEqual, 2)

		b.Put(6)
		b.Put(7)
		c.So(b.Size(), convey.ShouldEqual, 4)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, true)
		c.So(v, convey.ShouldEqual, 3)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, true)
		c.So(v, convey.ShouldEqual, 4)

		v, ok = b.Pop()
		c.So(ok, convey.ShouldEqual, true)
		c.So(v, convey.ShouldEqual, 6)
		c.So(b.Size(), convey.ShouldEqual, 1)
	})
}

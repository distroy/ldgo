/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAutoLocker(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		mu := &sync.Mutex{}

		l := AutoLock(mu)
		c.So(l.locked, convey.ShouldEqual, false)

		l.Unlock()
		c.So(l.locked, convey.ShouldEqual, false)

		c.So(l.TryLock(), convey.ShouldEqual, true)
		c.So(l.locked, convey.ShouldEqual, true)

		l.Lock()
		c.So(l.locked, convey.ShouldEqual, true)

		l.Unlock()
		c.So(l.locked, convey.ShouldEqual, false)

		l.Unlock()
		c.So(l.locked, convey.ShouldEqual, false)

		l.Lock()
		c.So(l.locked, convey.ShouldEqual, true)

		c.So(l.TryLock(), convey.ShouldEqual, true)
		c.So(l.locked, convey.ShouldEqual, true)

		l.Unlock()
		c.So(l.locked, convey.ShouldEqual, false)
	})
}

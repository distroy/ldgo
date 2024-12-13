/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGetPool(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		p0 := GetPool(func() []byte { return make([]byte, 0, 1024) })
		c.So(p0, convey.ShouldNotBeNil)
		p1 := GetPool(func() []byte { return make([]byte, 0, 2048) })
		c.So(p1, convey.ShouldNotBeNil)
		c.So(p0, convey.ShouldEqual, p1)
	})
}

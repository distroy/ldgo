/*
 * Copyright (C) distroy
 */

package ldmath

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(Sum(0), convey.ShouldEqual, 0)
		convey.So(Sum(123, -115, 35), convey.ShouldEqual, 123-115+35)
	})
}

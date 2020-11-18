/*
 * Copyright (C) distroy
 */

package ldgo

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_SortStrings(t *testing.T) {
	convey.Convey("", t, func() {
		l := []string{"223", "562", "424", "642", "223", "abc", "aab", "22", "cbd", "abc"}
		convey.So(IsSortedStrings(l), convey.ShouldBeFalse)

		SortStrings(l)

		convey.So(IsSortedStrings(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []string{
			"22", "223", "223", "424", "562", "642", "aab", "abc", "abc", "cbd",
		})

		convey.So(SearchStrings(l, ""), convey.ShouldEqual, 0)
		convey.So(SearchStrings(l, "123"), convey.ShouldEqual, 0)
		convey.So(SearchStrings(l, "24"), convey.ShouldEqual, 3)
		convey.So(SearchStrings(l, "zzz"), convey.ShouldEqual, 10)
	})
}

/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo/v3/ldmath"
	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/distroy/ldgo/v3/ldsort"
	"github.com/smartystreets/goconvey/convey"
)

func TestTopkStrings(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]string, 0, n)
			topk := make([]string, 0, k)

			for i := 0; i < n; i++ {
				x := ldrand.String(8)
				origin = append(origin, x)
				topk, _ = TopkStringsAdd(topk, k, x)
			}

			ldsort.SortStrings(origin)
			ldsort.SortStrings(topk)

			size := ldmath.MinInt(n, k)
			origin = origin[:size]
			convey.So(topk, convey.ShouldResemble, origin)
		})
	}
	convey.Convey(t.Name(), t, func() {
		testTopkFunc(10, 20)
		testTopkFunc(100, 5)
		testTopkFunc(100, 10)
	})
}

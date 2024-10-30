/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/distroy/ldgo/v3/ldsort"
	"github.com/smartystreets/goconvey/convey"
)

func testTopk(c convey.C, n, k int) {
	less := func(a, b int) bool { return a < b }

	name := fmt.Sprintf("topk-n:%d-k:%d", n, k)
	c.Convey(name, func(c convey.C) {
		origin := make([]int, 0, n)

		topk := New[int](k, less)

		for i := 0; i < n; i++ {
			x := ldrand.Intn(100)
			origin = append(origin, x)
			topk.Add(x)
		}

		ldsort.SortInts(origin)
		ldsort.SortInts(topk.Data())

		size := min(n, k)
		origin = origin[:size]
		c.So(topk.Data(), convey.ShouldResemble, origin)
	})
}

func TestTopK(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		testTopk(c, 10, 20)
		testTopk(c, 100, 5)
		testTopk(c, 100, 10)
	})
}

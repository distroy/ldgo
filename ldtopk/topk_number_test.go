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

func TestTopkInts(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]int, 0, n)
			topk := make([]int, 0, k)

			for i := 0; i < n; i++ {
				x := ldrand.Intn(100)
				origin = append(origin, x)
				topk, _ = TopkIntsAdd(topk, k, x)
			}

			ldsort.SortInts(origin)
			ldsort.SortInts(topk)

			size := min(n, k)
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

func TestTopkInt64s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]int64, 0, n)
			topk := make([]int64, 0, k)

			for i := 0; i < n; i++ {
				x := ldrand.Int63n(100)
				origin = append(origin, x)
				topk, _ = TopkInt64sAdd(topk, k, x)
			}

			ldsort.SortInt64s(origin)
			ldsort.SortInt64s(topk)

			size := min(n, k)
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

func TestTopkInt32s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]int32, 0, n)
			topk := make([]int32, 0, k)

			for i := 0; i < n; i++ {
				x := ldrand.Int31n(100)
				origin = append(origin, x)
				topk, _ = TopkInt32sAdd(topk, k, x)
			}

			ldsort.SortInt32s(origin)
			ldsort.SortInt32s(topk)

			size := min(n, k)
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

func TestTopkInt16s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]int16, 0, n)
			topk := make([]int16, 0, k)

			for i := 0; i < n; i++ {
				x := int16(ldrand.Int31n(100))
				origin = append(origin, x)
				topk, _ = TopkInt16sAdd(topk, k, x)
			}

			ldsort.SortInt16s(origin)
			ldsort.SortInt16s(topk)

			size := min(n, k)
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

func TestTopkInt8s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]int8, 0, n)
			topk := make([]int8, 0, k)

			for i := 0; i < n; i++ {
				x := int8(ldrand.Int31n(100))
				origin = append(origin, x)
				topk, _ = TopkInt8sAdd(topk, k, x)
			}

			ldsort.SortInt8s(origin)
			ldsort.SortInt8s(topk)

			size := min(n, k)
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

func TestTopkUints(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]uint, 0, n)
			topk := make([]uint, 0, k)

			for i := 0; i < n; i++ {
				x := uint(ldrand.Intn(100))
				origin = append(origin, x)
				topk, _ = TopkUintsAdd(topk, k, x)
			}

			ldsort.SortUints(origin)
			ldsort.SortUints(topk)

			size := min(n, k)
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

func TestTopkUint64s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]uint64, 0, n)
			topk := make([]uint64, 0, k)

			for i := 0; i < n; i++ {
				x := uint64(ldrand.Int63n(100))
				origin = append(origin, x)
				topk, _ = TopkUint64sAdd(topk, k, x)
			}

			ldsort.SortUint64s(origin)
			ldsort.SortUint64s(topk)

			size := min(n, k)
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

func TestTopkUint32s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]uint32, 0, n)
			topk := make([]uint32, 0, k)

			for i := 0; i < n; i++ {
				x := uint32(ldrand.Int63n(100))
				origin = append(origin, x)
				topk, _ = TopkUint32sAdd(topk, k, x)
			}

			ldsort.SortUint32s(origin)
			ldsort.SortUint32s(topk)

			size := min(n, k)
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

func TestTopkUint16s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]uint16, 0, n)
			topk := make([]uint16, 0, k)

			for i := 0; i < n; i++ {
				x := uint16(ldrand.Int63n(100))
				origin = append(origin, x)
				topk, _ = TopkUint16sAdd(topk, k, x)
			}

			ldsort.SortUint16s(origin)
			ldsort.SortUint16s(topk)

			size := min(n, k)
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

func TestTopkUint8s(t *testing.T) {
	testTopkFunc := func(n, k int) {
		name := fmt.Sprintf("%s-n:%d-k:%d", t.Name(), n, k)
		convey.Convey(name, func() {
			origin := make([]uint8, 0, n)
			topk := make([]uint8, 0, k)

			for i := 0; i < n; i++ {
				x := uint8(ldrand.Int63n(100))
				origin = append(origin, x)
				topk, _ = TopkUint8sAdd(topk, k, x)
			}

			ldsort.SortUint8s(origin)
			ldsort.SortUint8s(topk)

			size := min(n, k)
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

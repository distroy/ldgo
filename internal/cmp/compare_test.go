/*
 * Copyright (C) distroy
 */

package cmp

import (
	"math"
	"testing"
	"time"
	"unsafe"

	"github.com/smartystreets/goconvey/convey"
)

func TestCompareInterface(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("aType != bType", func(c convey.C) {
			c.So(Compare(nil, false), convey.ShouldEqual, -1)
			c.So(Compare(true, 0), convey.ShouldEqual, -1)

			c.So(Compare((*int)(nil), (*uint)(nil)), convey.ShouldEqual, -1)
			c.So(Compare((*int8)(nil), (*uint)(nil)), convey.ShouldEqual, -1)
			c.So(Compare((*int)(nil), (*int)(nil)), convey.ShouldEqual, 0)
			c.So(Compare(0, []int{}), convey.ShouldEqual, -1)
		})

		c.Convey("bool", func(c convey.C) {
			c.So(Compare(false, true), convey.ShouldEqual, -1)
			c.So(Compare(true, false), convey.ShouldEqual, 1)
			c.So(Compare(false, false), convey.ShouldEqual, 0)
			c.So(Compare(true, true), convey.ShouldEqual, 0)
		})

		c.Convey("int", func(c convey.C) {
			c.So(Compare(uint64(99), 100), convey.ShouldEqual, -1)
			c.So(Compare(uint64(math.MaxInt64+1), 100), convey.ShouldEqual, 1)
			c.So(Compare(200, int64(100)), convey.ShouldEqual, 1)
			c.So(Compare(200, uint64(200)), convey.ShouldEqual, 0)
			c.So(Compare(uint32(200), uint64(200)), convey.ShouldEqual, 0)
			c.So(Compare(int64(-200), 100), convey.ShouldEqual, -1)
			c.So(Compare(int64(-200), uint(100)), convey.ShouldEqual, -1)
		})

		c.Convey("float", func(c convey.C) {
			c.So(Compare(float64(100.0), float32(100.0)), convey.ShouldEqual, 0)
			c.So(Compare(99.0, 100.0), convey.ShouldEqual, -1)
			c.So(Compare(99.0, math.NaN()), convey.ShouldEqual, 1)
			c.So(Compare(-99.0, math.NaN()), convey.ShouldEqual, 1)
			c.So(Compare(-99.0, float32(math.NaN())), convey.ShouldEqual, 1)
			c.So(Compare(float32(math.NaN()), 100.0), convey.ShouldEqual, -1)

			c.So(Compare(math.NaN(), math.NaN()), convey.ShouldEqual, 0)
			c.So(Compare(float32(math.NaN()), float32(math.NaN())), convey.ShouldEqual, 0)

			c.So(Compare(-99.0, math.Inf(1)), convey.ShouldEqual, -1)
			c.So(Compare(-99.0, math.Inf(-1)), convey.ShouldEqual, 1)

			c.So(Compare(float64(4503599627370496.1), float64(4503599627370496)), convey.ShouldEqual, 0) // exceeds the precision of float64

			c.Convey("int", func(c convey.C) {
				c.So(Compare(float64(4503599627370496), int64(4503599627370496)), convey.ShouldEqual, 0)    //
				c.So(Compare(float64(4503599627370496.1), int64(4503599627370496)), convey.ShouldEqual, 0)  // exceeds the precision of float64
				c.So(Compare(float64(503599627370496.1), int64(503599627370496)), convey.ShouldEqual, 1)    //
				c.So(Compare(float64(45035996273704960), int64(45035996273704961)), convey.ShouldEqual, -1) //
				c.So(Compare(float64(45035996273704961), int64(45035996273704960)), convey.ShouldEqual, 0)  // exceeds the precision of float64
				c.So(Compare(float64(math.MaxInt64)*2, int64(45035996273704960)), convey.ShouldEqual, 1)    //

				c.So(Compare(float64(-4503599627370496), int64(-4503599627370496)), convey.ShouldEqual, 0)   //
				c.So(Compare(float64(-4503599627370496.1), int64(-4503599627370496)), convey.ShouldEqual, 0) // exceeds the precision of float64
				c.So(Compare(float64(-503599627370496.1), int64(-503599627370496)), convey.ShouldEqual, -1)  //
				c.So(Compare(float64(-45035996273704960), int64(-45035996273704961)), convey.ShouldEqual, 1) //
				c.So(Compare(float64(-45035996273704961), int64(-45035996273704960)), convey.ShouldEqual, 0) // exceeds the precision of float64
				c.So(Compare(float64(math.MinInt64)*2, int64(45035996273704960)), convey.ShouldEqual, -1)    //
			})

			c.Convey("uint", func(c convey.C) {
				c.So(Compare(float64(4503599627370496), uint64(4503599627370496)), convey.ShouldEqual, 0)    //
				c.So(Compare(float64(4503599627370496.1), uint64(4503599627370496)), convey.ShouldEqual, 0)  // exceeds the precision of float64
				c.So(Compare(float64(503599627370496.1), uint64(503599627370496)), convey.ShouldEqual, 1)    //
				c.So(Compare(float64(45035996273704960), uint64(45035996273704961)), convey.ShouldEqual, -1) //
				c.So(Compare(float64(45035996273704961), uint64(45035996273704960)), convey.ShouldEqual, 0)  // exceeds the precision of float64
				c.So(Compare(float64(math.MaxUint)*2, uint64(45035996273704960)), convey.ShouldEqual, 1)     //
				c.So(Compare(float64(-1), uint64(0)), convey.ShouldEqual, -1)                                //
			})
		})

		c.Convey("number", func(c convey.C) {
			c.So(Compare(99, float64(100)), convey.ShouldEqual, -1)
			c.So(Compare(99, float64(99)), convey.ShouldEqual, 0)
		})

		c.Convey("complex", func(c convey.C) {
			c.So(Compare(complex(100, 200), complex64(complex(100, 200))), convey.ShouldEqual, 0)

			c.So(Compare(complex(100, 200), complex(11, 300)), convey.ShouldEqual, 1)
			c.So(Compare(complex(100, 200), complex(111, -300)), convey.ShouldEqual, -1)
			c.So(Compare(complex(100, 200), complex(100, 300)), convey.ShouldEqual, -1)
			c.So(Compare(complex(100, 200), complex(100, 150)), convey.ShouldEqual, 1)
		})

		c.Convey("string", func(c convey.C) {
			c.So(Compare("", `abc`), convey.ShouldEqual, -1)
			c.So(Compare("aaa", `a`), convey.ShouldEqual, 1)
			c.So(Compare("bbb", `aaaaaa`), convey.ShouldEqual, 1)
		})

		c.Convey("map", func(c convey.C) {
			c.So(Compare(map[int]int{0: 0}, map[interface{}]int{0: 0}), convey.ShouldEqual, -1)
			c.So(Compare(map[int]int{0: 0}, map[int]int{0: 0}), convey.ShouldEqual, 0)
			c.So(Compare(map[int]int{0: 0}, map[int]int{}), convey.ShouldEqual, 1)
			c.So(Compare(map[int]int{0: 0}, map[int]int{1: 0}), convey.ShouldEqual, 1)
			c.So(Compare(map[int]int{1: 1}, map[int]int{1: 0}), convey.ShouldEqual, 1)
			c.So(Compare(map[int]int{1: 1}, map[int]int{0: 0, 1: 0}), convey.ShouldEqual, -1)
			c.So(Compare(map[int]int{1: 1}, map[int]int{1: 1, 2: 0}), convey.ShouldEqual, -1)
			c.So(Compare(map[int]int{0: 0, 1: 0}, map[int]int{0: 0}), convey.ShouldEqual, 1)
		})

		c.Convey("slice", func(c convey.C) {
			c.So(Compare(
				[]interface{}{100, uint(200), float32(300)},
				[]interface{}{100, float64(200), ""},
			), convey.ShouldEqual, -1)

			c.So(Compare(
				[]interface{}{100, uint(200), ""},
				[]interface{}{100, float64(200), ""},
			), convey.ShouldEqual, 0)
		})

		c.Convey("array", func(c convey.C) {
			c.So(Compare(
				[...]interface{}{100, uint(200), float32(300)},
				[...]interface{}{100, float64(200), ""},
			), convey.ShouldEqual, -1)

			c.So(Compare(
				[...]interface{}{100, uint(200), ""},
				[...]interface{}{100, float64(200), ""},
			), convey.ShouldEqual, 0)
		})

		c.Convey("pointer", func(c convey.C) {
			aa := 1
			bb := 2
			cc := 1
			c.So(Compare(&aa, &aa), convey.ShouldEqual, 0)
			c.So(Compare(&aa, &bb), convey.ShouldEqual, -1)
			c.So(Compare(&aa, &cc), convey.ShouldEqual, 0)
			c.So(Compare(&aa, (*int)(nil)), convey.ShouldEqual, 1)
		})
		c.Convey("unsafe pointer", func(c convey.C) {
			aa := 1
			bb := 2
			cc := 1
			c.So(Compare(unsafe.Pointer(&aa), unsafe.Pointer(&aa)), convey.ShouldEqual, 0)
			c.So(Compare(unsafe.Pointer(&aa), unsafe.Pointer(&bb)), convey.ShouldNotEqual, 0)
			c.So(Compare(unsafe.Pointer(&aa), unsafe.Pointer(&cc)), convey.ShouldNotEqual, 0)
			c.So(Compare(unsafe.Pointer(&aa), unsafe.Pointer(nil)), convey.ShouldEqual, 1)
		})

		c.Convey("chan", func(c convey.C) {
			aa := make(chan struct{})
			bb := make(chan struct{})
			c.So(Compare(aa, aa), convey.ShouldEqual, 0)
			c.So(Compare(aa, bb), convey.ShouldNotEqual, 0)
			c.So(Compare(aa, (chan struct{})(nil)), convey.ShouldEqual, 1)
		})
		c.Convey("func", func(c convey.C) {
			aa := func() {}
			bb := func() {}
			c.So(Compare(aa, aa), convey.ShouldEqual, 0)
			c.So(Compare(aa, bb), convey.ShouldNotEqual, 0)
			c.So(Compare(aa, (func())(nil)), convey.ShouldEqual, 1)
			c.So(Compare(aa, nil), convey.ShouldEqual, 1)
		})

		c.Convey("comparer", func(c convey.C) {
			c.So(Compare(time.Unix(0, 0), time.Unix(0, 0)), convey.ShouldEqual, 0)
			c.So(Compare(time.Unix(123, 0), time.Unix(123, 0)), convey.ShouldEqual, 0)
			c.So(Compare(time.Unix(123, 0).In(time.UTC), time.Unix(123, 0).In(time.Local)), convey.ShouldEqual, 0)
			c.So(Compare(time.Unix(0, 0), time.Unix(123, 0)), convey.ShouldEqual, -1)
			c.So(Compare(time.Unix(123, 0), time.Unix(0, 0)), convey.ShouldEqual, 1)
		})

		c.Convey("struct", func(c convey.C) {
			type StructA struct {
				Int int
			}
			type StructB struct {
				Int int
			}
			c.So(Compare(StructA{100}, StructB{100}), convey.ShouldEqual, -1)
			c.So(Compare(StructB{100}, StructA{100}), convey.ShouldEqual, +1)
			c.So(Compare(StructA{100}, StructA{100}), convey.ShouldEqual, +0)
			c.So(Compare(StructA{100}, StructA{99}), convey.ShouldEqual, +1)
		})
	})
}

func TestCompareBool(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareBool(false, false), convey.ShouldEqual, 0)
		c.So(CompareBool(true, true), convey.ShouldEqual, 0)
		c.So(CompareBool(false, true), convey.ShouldEqual, -1)
		c.So(CompareBool(true, false), convey.ShouldEqual, 1)
	})
}

func TestCompareByte(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[byte](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[byte](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[byte](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[byte](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareRune(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[rune](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[rune](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[rune](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[rune](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareInt(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable(int(0), 0), convey.ShouldEqual, 0)
		c.So(CompareComparable(int(123), 123), convey.ShouldEqual, 0)
		c.So(CompareComparable(int(0), 123), convey.ShouldEqual, -1)
		c.So(CompareComparable(int(123), 0), convey.ShouldEqual, 1)
	})
}

func TestCompareInt8(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[int8](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[int8](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[int8](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[int8](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareInt16(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[int16](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[int16](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[int16](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[int16](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareInt32(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[int32](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[int32](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[int32](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[int32](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareInt64(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[int64](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[int64](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[int64](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[int64](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareUint(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[uint](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[uint](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareUint8(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[uint8](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint8](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint8](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[uint8](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareUint16(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[uint16](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint16](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint16](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[uint16](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareUint32(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[uint32](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint32](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint32](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[uint32](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareUint64(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable[uint32](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint32](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[uint32](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[uint32](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareFloat32(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable(float32(math.NaN()), float32(math.NaN())), convey.ShouldEqual, 0)
		c.So(CompareComparable(0, float32(math.NaN())), convey.ShouldEqual, 1)
		c.So(CompareComparable(float32(math.NaN()), 0), convey.ShouldEqual, -1)

		c.So(CompareComparable(float32(math.Inf(-1)), float32(math.NaN())), convey.ShouldEqual, 1)
		c.So(CompareComparable(float32(math.NaN()), float32(math.Inf(-1))), convey.ShouldEqual, -1)

		c.So(CompareComparable(float32(math.Inf(1)), float32(math.Inf(1))), convey.ShouldEqual, 0)
		c.So(CompareComparable(float32(math.Inf(-1)), float32(math.Inf(-1))), convey.ShouldEqual, 0)

		c.So(CompareComparable(float32(math.Inf(-1)), float32(math.Inf(1))), convey.ShouldEqual, -1)
		c.So(CompareComparable(float32(math.Inf(1)), float32(math.Inf(-1))), convey.ShouldEqual, 1)

		c.So(CompareComparable(float32(math.Inf(-1)), 0), convey.ShouldEqual, -1)
		c.So(CompareComparable(0, float32(math.Inf(-1))), convey.ShouldEqual, 1)
		c.So(CompareComparable(float32(math.Inf(1)), 0), convey.ShouldEqual, 1)
		c.So(CompareComparable(0, float32(math.Inf(1))), convey.ShouldEqual, -1)

		c.So(CompareComparable[float32](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[float32](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[float32](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[float32](123, 0), convey.ShouldEqual, 1)
	})
}

func TestCompareFloat64(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparable(math.NaN(), math.NaN()), convey.ShouldEqual, 0)
		c.So(CompareComparable(0, math.NaN()), convey.ShouldEqual, 1)
		c.So(CompareComparable(math.NaN(), 0), convey.ShouldEqual, -1)

		c.So(CompareComparable(math.Inf(-1), math.NaN()), convey.ShouldEqual, 1)
		c.So(CompareComparable(math.NaN(), math.Inf(-1)), convey.ShouldEqual, -1)

		c.So(CompareComparable(math.Inf(1), math.Inf(1)), convey.ShouldEqual, 0)
		c.So(CompareComparable(math.Inf(-1), math.Inf(-1)), convey.ShouldEqual, 0)

		c.So(CompareComparable(math.Inf(-1), math.Inf(1)), convey.ShouldEqual, -1)
		c.So(CompareComparable(math.Inf(1), math.Inf(-1)), convey.ShouldEqual, 1)

		c.So(CompareComparable(math.Inf(-1), 0), convey.ShouldEqual, -1)
		c.So(CompareComparable(0, math.Inf(-1)), convey.ShouldEqual, 1)
		c.So(CompareComparable(math.Inf(1), 0), convey.ShouldEqual, 1)
		c.So(CompareComparable(0, math.Inf(1)), convey.ShouldEqual, -1)

		c.So(CompareComparable[float64](0, 0), convey.ShouldEqual, 0)
		c.So(CompareComparable[float64](123, 123), convey.ShouldEqual, 0)
		c.So(CompareComparable[float64](0, 123), convey.ShouldEqual, -1)
		c.So(CompareComparable[float64](123, 0), convey.ShouldEqual, 1)

		c.So(CompareComparable(float64(4503599627370496.1), 4503599627370496), convey.ShouldEqual, 0) // exceeds the precision of float64
	})
}

func TestString(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareString(``, ""), convey.ShouldEqual, 0)
		c.So(CompareString(`aaa`, "aaa"), convey.ShouldEqual, 0)

		c.So(CompareString(`aaa`, ""), convey.ShouldEqual, 1)
		c.So(CompareString(``, "aaa"), convey.ShouldEqual, -1)
		c.So(CompareString(`abc`, "aaa"), convey.ShouldEqual, 1)
		c.So(CompareString(`aa`, "aaa"), convey.ShouldEqual, -1)
	})
}

func TestCompareTime(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(CompareComparer(time.Unix(0, 0), time.Unix(0, 0)), convey.ShouldEqual, 0)
		c.So(CompareComparer(time.Unix(123, 0).In(time.UTC), time.Unix(123, 0).In(time.Local)), convey.ShouldEqual, 0)
		c.So(CompareComparer(time.Unix(123, 0), time.Unix(123, 0)), convey.ShouldEqual, 0)
		c.So(CompareComparer(time.Unix(0, 0), time.Unix(123, 0)), convey.ShouldEqual, -1)
		c.So(CompareComparer(time.Unix(123, 0), time.Unix(0, 0)), convey.ShouldEqual, 1)
	})
}

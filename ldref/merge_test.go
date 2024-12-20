/*
 * Copyright (C) distroy
 */

package ldref

import (
	"testing"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldref/internal/copybenchstruct1"
	"github.com/smartystreets/goconvey/convey"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/distroy/ldgo/v2/ldref
cpu: VirtualApple @ 2.50GHz
Benchmark_mergeV1
Benchmark_mergeV1-10                    21379281                54.55 ns/op
Benchmark_mergeV2
Benchmark_mergeV2-10                    23590855                54.51 ns/op
Benchmark_mergeV1WithClone
Benchmark_mergeV1WithClone-10              16126             98721 ns/op
Benchmark_mergeV2WithClone
Benchmark_mergeV2WithClone-10              21660             52806 ns/op
PASS
ok      github.com/distroy/ldgo/v2/ldref        14.325s
*/

type testErrorStruct struct {
	value interface{}
}

func (testErrorStruct) Error() string { return "" }

type testErrorStruct2 struct {
	value interface{}
}

func (*testErrorStruct2) Error() string { return "" }

func testMergeWithFunc(t *testing.T, fnMerge func(target, source interface{}, cfg ...*MergeConfig) error) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("fail", func(c convey.C) {
			c.Convey("to invalid type", func(c convey.C) {
				err := fnMerge(1, 2)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNotPtr)
			})
			c.Convey("to nil ptr", func(c convey.C) {
				err := fnMerge((*int)(nil), 2)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNilPtr)
			})
		})

		c.Convey("succ", func(c convey.C) {
			c.Convey("to interface", func(c convey.C) {
				c.Convey("from struct", func(c convey.C) {
					var target error
					source := testErrorStruct{value: "abcde"}

					err := fnMerge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, testErrorStruct{value: "abcde"})
				})
				c.Convey("from ptr to struct 1", func(c convey.C) {
					var target error
					source := &testErrorStruct{value: "abcde"}

					err := fnMerge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, testErrorStruct{value: "abcde"})
				})
				c.Convey("from ptr to struct 2", func(c convey.C) {
					var target error
					source := &testErrorStruct2{value: "abcde"}

					err := fnMerge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, &testErrorStruct2{value: "abcde"})
				})
			})

			c.Convey("from nil ptr", func(c convey.C) {
				var (
					target int = 1
				)

				err := fnMerge(&target, (*int)(nil))
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1)
			})

			c.Convey("normal type 1", func(c convey.C) {
				var (
					target int
					source = 1234
				)

				err := fnMerge(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1234)
			})
			c.Convey("normal type 2", func(c convey.C) {
				var (
					target int
					source = 1234
				)

				err := fnMerge(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1234)
			})

			c.Convey("ptr", func(c convey.C) {
				c.Convey("no clone", func(c convey.C) {
					var (
						target (*int)
						source = 1234
					)

					err := fnMerge(&target, &source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, &source)
				})
				c.Convey("clone", func(c convey.C) {
					var (
						target (*int)
						source = 1234
					)

					err := fnMerge(&target, &source, &MergeConfig{Clone: true})
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldNotEqual, &source)
					c.So(target, convey.ShouldResemble, &source)
				})
			})

			c.Convey("struct", func(c convey.C) {
				var (
					target = &testCloneStruct{
						String: "abc",
					}
					source = &testCloneStruct{
						Int:    1234,
						String: "xyz",
					}
				)

				err := fnMerge(target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, &testCloneStruct{
					Int:    1234,
					String: "abc",
				})
			})

			c.Convey("map", func(c convey.C) {
				var (
					target = map[string]any{
						"a": 1234,
						"c": &testCloneStruct{
							String: "abc",
						},
					}
					source = map[string]any{
						"a": 2345,
						"b": "abc",
						"c": &testCloneStruct{
							Int:    1234,
							String: "xyz",
						},
					}
				)

				err := fnMerge(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, map[string]any{
					"a": 1234,
					"b": "abc",
					"c": &testCloneStruct{
						Int:    1234,
						String: "abc",
					},
				})
			})

			c.Convey("slice", func(c convey.C) {
				c.Convey("no merge elem", func(c convey.C) {
					var (
						target = map[string]any{
							"a": []any(nil),
						}
						source = map[string]any{
							"a": []any{1, 2, 4},
						}
					)

					err := fnMerge(&target, &source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, map[string]any{
						"a": []any{1, 2, 4},
					})
				})

				c.Convey("merge elem", func(c convey.C) {
					var (
						target = map[string]any{
							"a": []any{0, 3, 0},
						}
						source = map[string]any{
							"a": []any{1, 2, 4, 7},
						}
					)

					err := fnMerge(&target, &source, &MergeConfig{MergeSlice: true})
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, map[string]any{
						"a": []any{1, 3, 4, 7},
					})
				})
			})

			c.Convey("array", func(c convey.C) {
				c.Convey("no merge elem", func(c convey.C) {
					var (
						target = [4]any{}
						source = [4]any{1, 2, 4}
					)

					err := fnMerge(&target, &source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, [4]any{1, 2, 4})
				})
				c.Convey("merge elem", func(c convey.C) {
					var (
						target = [4]any{0, 0, 5}
						source = [4]any{0, 2, 0, 14}
					)

					err := fnMerge(&target, &source, &MergeConfig{MergeArray: true})
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, [4]any{0, 2, 5, 14})
				})
			})
		})
	})
}

func TestMerge(t *testing.T)    { testMergeWithFunc(t, Merge) }
func Test_mergeV1(t *testing.T) { testMergeWithFunc(t, mergeV1) }
func Test_mergeV2(t *testing.T) { testMergeWithFunc(t, mergeV2) }

func benchMergeFunc(b *testing.B, fnMerge func(target, source interface{}, cfg ...*MergeConfig) error, clone bool) {
	size := 1024
	mask := size - 1
	objs := benchPrepareObjects(size)
	{
		target := &copybenchstruct1.ItemCardData{}
		source := objs[0]
		fnMerge(target, source)
	}
	cfg := &MergeConfig{
		Clone: clone,
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		count := 0
		for p.Next() {
			idx := count & mask
			count++

			target := &copybenchstruct1.ItemCardData{}
			source := objs[idx]
			fnMerge(target, source, cfg)
		}
	})
	b.StopTimer()
}

func Benchmark_mergeV1(b *testing.B)          { benchMergeFunc(b, mergeV1, false) }
func Benchmark_mergeV2(b *testing.B)          { benchMergeFunc(b, mergeV2, false) }
func Benchmark_mergeV1WithClone(b *testing.B) { benchMergeFunc(b, mergeV1, true) }
func Benchmark_mergeV2WithClone(b *testing.B) { benchMergeFunc(b, mergeV2, true) }

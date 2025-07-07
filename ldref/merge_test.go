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
BenchmarkMerge/v1-empty-target-14               40739823                29.41 ns/op
BenchmarkMerge/v2-empty-target-14               46734962                29.65 ns/op
BenchmarkMerge/v1-value-target-14                 113146             10999 ns/op
BenchmarkMerge/v2-value-target-14                 208208              6498 ns/op
BenchmarkMerge/v1-empty-target-with-clone-14               40005             29611 ns/op
BenchmarkMerge/v2-empty-target-with-clone-14               40640             29867 ns/op
BenchmarkMerge/v1-value-target-with-clone-14               68564             17961 ns/op
BenchmarkMerge/v2-value-target-with-clone-14              212548              6226 ns/op
PASS
ok      github.com/distroy/ldgo/v2/ldref        18.859s
*/

type testErrorStruct struct {
	value any
}

func (testErrorStruct) Error() string { return "" }

type testErrorStruct2 struct {
	value any
}

func (*testErrorStruct2) Error() string { return "" }

func testMergeWithFunc(t *testing.T, fnMerge func(target, source any, cfg ...*MergeConfig) error) {
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
						Boolp:  testNew(false),
					}
					source = &testCloneStruct{
						String: "xyz",
						Int:    1234,
						Uintp:  testNew[uint](2345),
						Boolp:  testNew(true),
					}
				)

				err := fnMerge(target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, &testCloneStruct{
					String: "abc",
					Int:    1234,
					Uintp:  testNew[uint](2345),
					Boolp:  testNew(false),
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

type benchMergeFuncConfig struct {
	merge       func(target, source any, cfg ...*MergeConfig) error
	emptyTarget bool
	config      *MergeConfig
}

func benchMergeFunc(b *testing.B, c *benchMergeFuncConfig) {
	var (
		merge = c.merge
		cfg   = c.config
	)
	size := 1024
	mask := size - 1
	objs := benchPrepareObjects(size)
	{
		target := objs[0]
		source := objs[1]
		merge(target, source)
	}
	getTarget := func() *copybenchstruct1.ItemCardData { return &copybenchstruct1.ItemCardData{} }
	if !c.emptyTarget {
		target := DeepClone(objs[0])
		getTarget = func() *copybenchstruct1.ItemCardData { return target }
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		count := 0
		for p.Next() {
			idx := count & mask
			count++

			target := getTarget()
			source := objs[idx]
			merge(target, source, cfg)
		}
	})
	b.StopTimer()
}

func BenchmarkMerge(b *testing.B) {
	cfg := func(merge func(target, source any, cfg ...*MergeConfig) error, clone bool, emptyTarget bool) *benchMergeFuncConfig {
		return &benchMergeFuncConfig{
			merge:       merge,
			emptyTarget: emptyTarget,
			config: &MergeConfig{
				Clone:      clone,
				MergeArray: true,
				MergeSlice: true,
			},
		}
	}
	b.Run("v1-empty-target", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV1, false, true)) })
	b.Run("v2-empty-target", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV2, false, true)) })
	b.Run("v1-value-target", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV1, false, false)) })
	b.Run("v2-value-target", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV2, false, false)) })
	b.Run("v1-empty-target-with-clone", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV1, true, true)) })
	b.Run("v2-empty-target-with-clone", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV2, true, true)) })
	b.Run("v1-value-target-with-clone", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV1, true, false)) })
	b.Run("v2-value-target-with-clone", func(b *testing.B) { benchMergeFunc(b, cfg(mergeV2, true, false)) })
}

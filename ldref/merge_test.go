/*
 * Copyright (C) distroy
 */

package ldref

import (
	"testing"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/smartystreets/goconvey/convey"
)

type testErrorStruct struct {
	value interface{}
}

func (testErrorStruct) Error() string { return "" }

type testErrorStruct2 struct {
	value interface{}
}

func (*testErrorStruct2) Error() string { return "" }

func TestMerge(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("fail", func(c convey.C) {
			c.Convey("to invalid type", func(c convey.C) {
				err := Merge(1, 2)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNotPtr)
			})
			c.Convey("to nil ptr", func(c convey.C) {
				err := Merge((*int)(nil), 2)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNilPtr)
			})
		})

		c.Convey("succ", func(c convey.C) {
			c.Convey("to interface", func(c convey.C) {
				c.Convey("from struct", func(c convey.C) {
					var target error
					source := testErrorStruct{value: "abcde"}

					err := Merge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, testErrorStruct{value: "abcde"})
				})
				c.Convey("from ptr to struct 1", func(c convey.C) {
					var target error
					source := &testErrorStruct{value: "abcde"}

					err := Merge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, testErrorStruct{value: "abcde"})
				})
				c.Convey("from ptr to struct 2", func(c convey.C) {
					var target error
					source := &testErrorStruct2{value: "abcde"}

					err := Merge(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, &testErrorStruct2{value: "abcde"})
				})
			})

			c.Convey("from nil ptr", func(c convey.C) {
				var (
					target int = 1
				)

				err := Merge(&target, (*int)(nil))
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1)
			})

			c.Convey("normal type 1", func(c convey.C) {
				var (
					target int
					source = 1234
				)

				err := Merge(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1234)
			})
			c.Convey("normal type 2", func(c convey.C) {
				var (
					target int
					source = 1234
				)

				err := Merge(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 1234)
			})

			c.Convey("ptr", func(c convey.C) {
				c.Convey("no clone", func(c convey.C) {
					var (
						target (*int)
						source = 1234
					)

					err := Merge(&target, &source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, &source)
				})
				c.Convey("clone", func(c convey.C) {
					var (
						target (*int)
						source = 1234
					)

					err := Merge(&target, &source, &MergeConfig{Clone: true})
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

				err := Merge(target, source)
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

				err := Merge(&target, &source)
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

					err := Merge(&target, &source)
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

					err := Merge(&target, &source, &MergeConfig{MergeSlice: true})
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

					err := Merge(&target, &source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, [4]any{1, 2, 4})
				})
				c.Convey("merge elem", func(c convey.C) {
					var (
						target = [4]any{0, 0, 5}
						source = [4]any{0, 2, 0, 14}
					)

					err := Merge(&target, &source, &MergeConfig{MergeArray: true})
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldResemble, [4]any{0, 2, 5, 14})
				})
			})
		})
	})
}

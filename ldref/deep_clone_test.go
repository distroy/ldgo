/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"sync"
	"testing"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/smartystreets/goconvey/convey"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/distroy/ldgo/v2/ldref
cpu: VirtualApple @ 2.50GHz
Benchmark_deepCloneV1
Benchmark_deepCloneV1-10           15806             75251 ns/op
Benchmark_deepCloneV2
Benchmark_deepCloneV2-10           22030             55012 ns/op
PASS
ok      github.com/distroy/ldgo/v2/ldref        6.935s
*/

func testDeepClone(t *testing.T, cloneFunc func(x0 any) any) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("nil", func(c convey.C) {
			c.Convey("interface{}", func(c convey.C) {
				v0 := interface{}(nil)
				// v1 := DeepClone(v0)
				v1 := cloneWithFuncForTest(v0, cloneFunc)
				c.So(v1, convey.ShouldBeNil)
			})
			c.Convey("error", func(c convey.C) {
				v0 := error(nil)
				// v1 := DeepClone(v0)
				v1 := cloneWithFuncForTest(v0, cloneFunc)
				c.So(v1, convey.ShouldBeNil)
			})
		})
		c.Convey("reflect.Value(*int)", func(c convey.C) {
			p0 := new(int)
			*p0 = 12345
			v0 := reflect.ValueOf(p0)
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1.Interface(), convey.ShouldResemble, v0.Interface())
			c.So(v1.Interface(), convey.ShouldResemble, v0.Interface())
		})
		c.Convey("*int", func(c convey.C) {
			v0 := new(int)
			*v0 = 12345
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("[]*int", func(c convey.C) {
			ptrToInt := func(n int) *int { return &n }
			v0 := []*int{ptrToInt(1), ptrToInt(2), ptrToInt(3)}
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
			d1 := v1
			for i := range v0 {
				c.So(d1[i], convey.ShouldNotEqual, v0[i])
			}
		})
		c.Convey("[3]int", func(c convey.C) {
			v0 := [3]int{1, 2, 3}
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("map[string]*int", func(c convey.C) {
			ptrToInt := func(n int) *int { return &n }
			v0 := map[string]*int{
				"a": ptrToInt(1),
				"b": ptrToInt(2),
				"c": ptrToInt(3),
				// "d": nil,
			}
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
			d1 := v1
			for i := range v0 {
				if d1[i] == nil {
					continue
				}
				c.So(d1[i], convey.ShouldNotEqual, v0[i])
				c.So(d1[i], convey.ShouldResemble, v0[i])
			}
		})
		c.Convey("struct", func(c convey.C) {
			v0 := testCloneStruct{
				String: "abc",
				Int:    123,
				Struct: &testCloneStruct{
					String: "xyz",
					Int:    234,
				},
			}
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
			d1 := v1
			c.So(d1.Struct, convey.ShouldNotEqual, v0.Struct)
		})
		c.Convey("*struct", func(c convey.C) {
			v0 := &testCloneStruct{
				String: "abc",
				Int:    123,
				Struct: &testCloneStruct{
					String: "xyz",
					Int:    234,
				},
			}
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
			d1 := v1
			c.So(d1.Struct, convey.ShouldNotEqual, v0.Struct)
		})

		c.Convey("error", func(c convey.C) {
			err := lderr.ErrUnkown
			v0 := lderr.New(err.Status(), err.Code(), err.Error())
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})

		c.Convey("*sync.Mutex", func(c convey.C) {
			v0 := &sync.Mutex{}
			v0.Lock()
			// v1 := DeepClone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldNotResemble, v0)
			c.So(v1, convey.ShouldResemble, &sync.Mutex{})
		})
	})
}
func TestDeepClone(t *testing.T)    { testDeepClone(t, DeepClone[any]) }
func Test_deepCloneV1(t *testing.T) { testDeepClone(t, deepCloneV1[any]) }
func Test_deepCloneV2(t *testing.T) { testDeepClone(t, deepCloneV2[any]) }

func Benchmark_deepCloneV1(b *testing.B) { benchCloneFunc(b, deepCloneV1[any]) }
func Benchmark_deepCloneV2(b *testing.B) { benchCloneFunc(b, deepCloneV2[any]) }

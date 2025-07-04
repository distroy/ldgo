/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"testing"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/smartystreets/goconvey/convey"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/distroy/ldgo/v2/ldref
cpu: VirtualApple @ 2.50GHz
Benchmark_cloneV1
Benchmark_cloneV1-10            47517410                22.76 ns/op
Benchmark_cloneV2
Benchmark_cloneV2-10            33064568                39.95 ns/op
PASS
ok      github.com/distroy/ldgo/v2/ldref        7.611s
*/

type testCloneStruct struct {
	String string
	Int    int
	Boolp  *bool
	Struct *testCloneStruct
}

func cloneWithFuncForTest[T any](v T, cloneFunc func(v any) any) T {
	vv := cloneFunc(v)
	if vv == nil {
		var x T
		return x
	}
	return vv.(T)
}

func testCloneFunc(t *testing.T, cloneFunc func(v any) any) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("nil", func(c convey.C) {
			c.Convey("interface{}", func(c convey.C) {
				v0 := any(nil)
				// v1 := Clone(v0)
				v1 := cloneWithFuncForTest(v0, cloneFunc)
				c.So(v1, convey.ShouldBeNil)
			})
			c.Convey("error", func(c convey.C) {
				v0 := error(nil)
				// v1 := Clone(v0)
				v1 := cloneWithFuncForTest(v0, cloneFunc)
				c.So(v1, convey.ShouldBeNil)
			})
		})
		c.Convey("reflect.Value(*int)", func(c convey.C) {
			p0 := new(int)
			*p0 = 12345
			v0 := reflect.ValueOf(p0)
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1.Interface(), convey.ShouldResemble, v0.Interface())
			c.So(v1.Interface(), convey.ShouldResemble, v0.Interface())
		})
		c.Convey("*int", func(c convey.C) {
			v0 := new(int)
			*v0 = 12345
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("[]int", func(c convey.C) {
			v0 := []int{1, 2, 3}
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("[3]int", func(c convey.C) {
			v0 := [3]int{1, 2, 3}
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
			c.So(&v1[0], convey.ShouldNotEqual, &v0[0])
		})
		c.Convey("map[string]int", func(c convey.C) {
			v0 := map[string]int{
				"a": 1,
				"b": 2,
			}
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("struct", func(c convey.C) {
			v0 := testCloneStruct{
				String: "abc",
				Int:    123,
			}
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("*struct", func(c convey.C) {
			v0 := &testCloneStruct{
				String: "abc",
				Int:    123,
			}
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
		c.Convey("error", func(c convey.C) {
			err := lderr.ErrUnkown
			v0 := lderr.New(err.Status(), err.Code(), err.Error())
			// v1 := Clone(v0)
			v1 := cloneWithFuncForTest(v0, cloneFunc)
			c.So(v1, convey.ShouldNotEqual, v0)
			c.So(v1, convey.ShouldResemble, v0)
		})
	})
}

func TestClone(t *testing.T)    { testCloneFunc(t, Clone[any]) }
func Test_cloneV1(t *testing.T) { testCloneFunc(t, cloneV1[any]) }
func Test_cloneV2(t *testing.T) { testCloneFunc(t, cloneV2[any]) }

func benchCloneFunc(b *testing.B, cloneFunc func(v any) any) {
	size := 1024
	mask := size - 1
	objs := benchPrepareObjects(size)
	{
		x0 := objs[0]
		cloneWithFuncForTest(x0, cloneFunc)
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		count := 0
		for p.Next() {
			var (
				idx = count
				obj = objs[idx&mask]
			)
			count++
			cloneWithFuncForTest(obj, cloneFunc)
		}
	})
	b.StopTimer()
}

func Benchmark_cloneV1(b *testing.B) { benchCloneFunc(b, cloneV1[any]) }
func Benchmark_cloneV2(b *testing.B) { benchCloneFunc(b, cloneV2[any]) }

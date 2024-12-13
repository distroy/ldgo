/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"sort"
	"testing"
	"unsafe"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/smartystreets/goconvey/convey"
)

func testDeepCopyFunc(t *testing.T, copyFunc func(target, source interface{}, cfg ...*CopyConfig) error) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("int to int", func(c convey.C) {
			var (
				target int
				source int = 100
			)

			err := copyFunc(target, source)
			c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNotPtr)
			c.So(target, convey.ShouldEqual, 0)
		})

		c.Convey("to reflect.Value", func(c convey.C) {
			c.Convey("int to reflect.Value - int (can not addr)", func(c convey.C) {
				var (
					target reflect.Value = reflect.ValueOf(int(0))
					source int           = 100
				)

				err := copyFunc(target, source)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNotPtr)
				c.So(target.Interface(), convey.ShouldEqual, 0)
			})

			c.Convey("int to reflect.Value - *int", func(c convey.C) {
				var (
					xxxx   *int          = new(int)
					target reflect.Value = reflect.ValueOf(xxxx)
					source int           = 100
				)

				err := copyFunc(target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target.Interface(), convey.ShouldResemble, &source)
			})

			c.Convey("int to reflect.Value - *int (can not addr)", func(c convey.C) {
				var (
					xxxx   *int
					target reflect.Value = reflect.ValueOf(xxxx)
					source int           = 100
				)

				err := copyFunc(target, source)
				// c.So(err, convey.ShouldBeNil)
				// c.So(target.Interface(), convey.ShouldEqual, 100)
				c.So(err, convey.ShouldResemble, lderr.ErrReflectTargetNilPtr)
				c.So(target.Interface(), convey.ShouldBeNil)
			})

			c.Convey("int to reflect.Value - *int (elem of **int)", func(c convey.C) {
				var (
					xxxx   **int         = new(*int)
					target reflect.Value = reflect.ValueOf(xxxx).Elem()
					source int           = 100
				)

				err := copyFunc(target, source)
				// c.So(err, convey.ShouldBeNil)
				// c.So(target.Interface(), convey.ShouldEqual, 100)
				c.So(err, convey.ShouldBeNil)
				c.So(target.Interface(), convey.ShouldResemble, &source)
			})

			c.Convey("int to reflect.Value - **int", func(c convey.C) {
				var (
					xxxx   *int
					target reflect.Value = reflect.ValueOf(&xxxx)
					source int           = 100
				)

				err := copyFunc(target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(xxxx, convey.ShouldResemble, &source)
			})
		})

		c.Convey("to *interface{}", func(c convey.C) {
			c.Convey("nil to *interface{}", func(c convey.C) {
				var (
					target interface{}
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, nil)
			})

			c.Convey("int to *interface{}", func(c convey.C) {
				var (
					target interface{}
					source int = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source)
			})

			c.Convey("*int to *interface{}", func(c convey.C) {
				var (
					target interface{}
					source *int = new(int)
				)
				*source = 100

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldNotEqual, source)
				c.So(target, convey.ShouldResemble, source)
			})

			c.Convey("**int to *interface{}", func(c convey.C) {
				var (
					target interface{}
					xxxx0  int   = 100
					xxxx1  *int  = &xxxx0
					source **int = &xxxx1
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldNotEqual, source)
				c.So(target, convey.ShouldResemble, source)
			})
		})

		c.Convey("to **int", func(c convey.C) {
			c.Convey("bool to **int", func(c convey.C) {
				var target *int = new(int)
				var source bool = true

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(*target, convey.ShouldEqual, 1)
			})

			c.Convey("**float to **int", func(c convey.C) {
				var (
					target *int     = new(int)
					source *float64 = new(float64)
				)
				*source = 100

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(*target, convey.ShouldEqual, *source)
			})

			c.Convey("**int to **int", func(c convey.C) {
				var (
					target *int = new(int)
					source *int = new(int)
				)
				*source = 100

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, source)
				c.So(target, convey.ShouldNotEqual, source)
			})

			c.Convey("nil to **int", func(c convey.C) {
				var (
					target *int
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldBeNil)
			})
		})

		c.Convey("func", func(c convey.C) {
			c.Convey("func to *func", func(c convey.C) {
				var (
					target func()
					source func() = func() {}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source)
			})

			c.Convey("func to *int", func(c convey.C) {
				var (
					target int    = 100
					source func() = func() {}
				)

				err := copyFunc(&target, source)
				c.So(err.Error(), convey.ShouldEqual, "func() can not copy to *int")
				c.So(target, convey.ShouldEqual, 0)
			})

			c.Convey("func to *string", func(c convey.C) {
				var (
					target string
					source func(interface{}) bool = IsZero
				)

				err := copyFunc(&target, source)
				c.So(err.Error(), convey.ShouldEqual, "func(interface {}) bool can not copy to *string")
				// c.So(err, convey.ShouldBeNil)
				// c.So(target, convey.ShouldEqual, "github.com/distroy/ldgo/v2/ldref.IsZero")
			})

			c.Convey("func to unsafe.Pointer", func(c convey.C) {
				var (
					target unsafe.Pointer
					source func(interface{}) bool = IsZero
				)

				err := copyFunc(&target, source)
				c.So(err.Error(), convey.ShouldEqual, "func(interface {}) bool can not copy to *unsafe.Pointer")
			})
		})

		c.Convey("to *bool", func(c convey.C) {
			c.Convey("nil to *bool", func(c convey.C) {
				var (
					target bool
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, false)
			})

			c.Convey("float to *bool", func(c convey.C) {
				var (
					target bool
					source float64 = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source != 0)
			})
			c.Convey("int to *bool", func(c convey.C) {
				var (
					target bool
					source int = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source != 0)
			})
			c.Convey("uint to *bool", func(c convey.C) {
				var (
					target bool
					source uint = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source != 0)
			})
			c.Convey("complex to *bool", func(c convey.C) {
				var (
					target bool
					source complex128 = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, source != 0)
			})
		})

		c.Convey("number to *numbers", func(c convey.C) {
			c.Convey("to *int", func(c convey.C) {
				c.Convey("nil to *int", func(c convey.C) {
					var (
						target int
						source interface{}
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, 0)
				})

				c.Convey("int to *int", func(c convey.C) {
					var (
						target int
						source int = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("uint to *int", func(c convey.C) {
					var (
						target int
						source uint = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("float64 to *int", func(c convey.C) {
					var (
						target int
						source float64 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("complex128 to *int", func(c convey.C) {
					var (
						target int
						source complex128 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})
			})
			c.Convey("to *uint", func(c convey.C) {
				c.Convey("nil to *uint", func(c convey.C) {
					var (
						target uint
						source interface{}
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, 0)
				})

				c.Convey("int to *uint", func(c convey.C) {
					var (
						target uint
						source int = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("uint to *uint", func(c convey.C) {
					var (
						target uint
						source uint = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("float64 to *uint", func(c convey.C) {
					var (
						target uint
						source float64 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("complex128 to *uint", func(c convey.C) {
					var (
						target uint
						source complex128 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})
			})
			c.Convey("to *float64", func(c convey.C) {
				c.Convey("nil to *float64", func(c convey.C) {
					var (
						target float64
						source interface{}
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, 0)
				})

				c.Convey("int to *float64", func(c convey.C) {
					var (
						target float64
						source int = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("uint to *float64", func(c convey.C) {
					var (
						target float64
						source uint = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("float64 to *float64", func(c convey.C) {
					var (
						target float64
						source float64 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("complex128 to *float64", func(c convey.C) {
					var (
						target float64
						source complex128 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})
			})
			c.Convey("to *complex128", func(c convey.C) {
				c.Convey("nil to *complex128", func(c convey.C) {
					var (
						target complex128
						source interface{}
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, 0)
				})

				c.Convey("int to *complex128", func(c convey.C) {
					var (
						target complex128
						source int = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("uint to *complex128", func(c convey.C) {
					var (
						target complex128
						source uint = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("float64 to *complex128", func(c convey.C) {
					var (
						target complex128
						source float64 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})

				c.Convey("complex128 to *complex128", func(c convey.C) {
					var (
						target complex128
						source complex128 = 100
					)

					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(target, convey.ShouldEqual, source)
				})
			})
		})

		c.Convey("number to *string", func(c convey.C) {
			c.Convey("nil to *string", func(c convey.C) {
				var (
					target string
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "")
			})

			c.Convey("bool to *string", func(c convey.C) {
				var (
					target string
					source bool = true
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "true")
			})

			c.Convey("int to *string", func(c convey.C) {
				var (
					target string
					source int = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "100")
			})
			c.Convey("uint to *string", func(c convey.C) {
				var (
					target string
					source uint = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "100")
			})
			c.Convey("float64 to *string", func(c convey.C) {
				var (
					target string
					source float64 = 100
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "100")
			})

			// c.Convey("complex128 to *string", func (c convey.C) {
			// 	var (
			// 		target string
			// 		source complex128 = 100
			// 	)
			//
			// 	err := copyFunc(&target, source)
			// 	c.So(err, convey.ShouldBeNil)
			// 	c.So(target, convey.ShouldEqual, "(100+0i)")
			// })
		})

		c.Convey("string to *number", func(c convey.C) {
			c.Convey("string to *float", func(c convey.C) {
				var (
					target float64
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 100)
			})
			c.Convey("string to *bool", func(c convey.C) {
				var (
					target bool
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, true)
			})
			c.Convey("string to *int", func(c convey.C) {
				var (
					target int
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 100)
			})
			c.Convey("string to *uint", func(c convey.C) {
				var (
					target uint
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, 100)
			})

			// c.Convey("string to *complex128", func (c convey.C) {
			// 	var (
			// 		target complex128
			// 		source string = "100"
			// 	)
			//
			// 	err := copyFunc(&target, source)
			// 	c.So(err, convey.ShouldBeNil)
			// 	c.So(target, convey.ShouldEqual, complex128(100))
			// })
		})

		c.Convey("string/[]byte/[]rune", func(c convey.C) {
			c.Convey("[]byte to *string", func(c convey.C) {
				var (
					target string
					source []byte = []byte("100")
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "100")
			})

			c.Convey("[]rune to *string", func(c convey.C) {
				var (
					target string
					source []rune = []rune("100")
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldEqual, "100")
			})

			c.Convey("string to *[]byte", func(c convey.C) {
				var (
					target []byte
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, []byte(source))
			})

			c.Convey("string to *[]rune", func(c convey.C) {
				var (
					target []rune
					source string = "100"
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, []rune(source))
			})
		})

		c.Convey("to *map", func(c convey.C) {
			c.Convey("nil to *map[int64]string", func(c convey.C) {
				var (
					target map[int64]string
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldBeNil)
			})

			c.Convey("map[int]string to *map[interface{}]interface{}", func(c convey.C) {
				var (
					target map[interface{}]interface{}
					source map[int]string = map[int]string{
						1: "abc",
						2: "xyz",
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, map[interface{}]interface{}{
					1: "abc",
					2: "xyz",
				})
			})

			c.Convey("map[int]string to *map[int64][]byte", func(c convey.C) {
				var (
					target map[int64][]byte
					source map[int]string = map[int]string{
						1: "abc",
						2: "xyz",
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, map[int64][]byte{
					1: []byte("abc"),
					2: []byte("xyz"),
				})
			})

			c.Convey("*struct to *map[string]string", func(c convey.C) {
				var (
					target map[string]string
					source *testCopyStruct = &testCopyStruct{
						Id:     123,
						Ignore: true,
						Name:   "abc",
						Age:    100,
						Female: true,
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, map[string]string{
					`id`:     "123",
					`Name`:   "abc",
					`Age`:    "100",
					`Female`: "true",
				})
			})

			c.Convey("[]string to *map[string]string", func(c convey.C) {
				var (
					target map[string]string
					source []string = []string{
						"abc", "xyz", "123", "zzz",
					}
				)

				err := copyFunc(&target, source)
				c.So(err.Error(), convey.ShouldEqual, "[]string can not copy to *map[string]string")
				// c.So(err, convey.ShouldBeNil)
				// c.So(target, convey.ShouldResemble, map[string]string{
				// 	`0`: "abc",
				// 	`1`: "xyz",
				// 	`2`: "123",
				// 	`3`: "zzz",
				// })
			})

			c.Convey("[]string to *map[string]struct{}", func(c convey.C) {
				var (
					target map[string]struct{}
					source []string = []string{
						"abc", "xyz", "123", "zzz",
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, map[string]struct{}{
					"abc": {},
					"xyz": {},
					"123": {},
					"zzz": {},
				})
			})
		})

		c.Convey("to *slice", func(c convey.C) {
			c.Convey("nil to *[]string", func(c convey.C) {
				var (
					target []string
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldBeNil)
			})

			c.Convey("nil to *[2]string", func(c convey.C) {
				var (
					target [2]string
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, [2]string{})
			})

			c.Convey("[]int64 to *[]complex128", func(c convey.C) {
				var (
					target []complex128
					source []int64 = []int64{
						2, -100, 356,
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, []complex128{
					2, -100, 356,
				})
			})

			c.Convey("map[string]struct{} to *[]string", func(c convey.C) {
				var (
					target []string
					source map[string]struct{} = map[string]struct{}{
						"abc": {},
						"xyz": {},
						"AAA": {},
						"123": {},
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)

				sort.Strings(target)
				c.So(target, convey.ShouldResemble, []string{
					"123", "AAA", "abc", "xyz",
				})
			})

		})

		c.Convey("to *struct/**struct", func(c convey.C) {
			c.Convey("nil to *struct", func(c convey.C) {
				var (
					target testCopyStruct
					source interface{}
				)

				err := copyFunc(&target, source)
				c.So(err.Error(), convey.ShouldEqual, "nil can not copy to *ldref.testCopyStruct")
				// c.So(err, convey.ShouldBeNil)
				// c.So(&target, convey.ShouldResemble, &testCopyStruct{})
			})

			c.Convey("struct to *struct", func(c convey.C) {
				var (
					target testCopyStruct
					source = testCopyStruct{
						Id:        100,
						Ignore:    true,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(&target, convey.ShouldResemble, &testCopyStruct{
					Id:        100,
					Name:      "abc",
					Age:       23,
					Female:    true,
					unexpored: new(int),
				})
			})

			c.Convey("*struct to *struct", func(c convey.C) {
				var (
					target testCopyStruct
					source *testCopyStruct = &testCopyStruct{
						Id:        100,
						Ignore:    true,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					}
				)

				err := copyFunc(&target, source)
				c.So(err, convey.ShouldBeNil)
				c.So(&target, convey.ShouldResemble, &testCopyStruct{
					Id:        100,
					Name:      "abc",
					Age:       23,
					Female:    true,
					unexpored: new(int),
				})
			})

			c.Convey("**struct to *struct", func(c convey.C) {
				var (
					target testCopyStruct
					source *testCopyStruct = &testCopyStruct{
						Id:     100,
						Ignore: true,
						Name:   "abc",
						Age:    23,
						Female: true,
					}
				)

				err := copyFunc(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(&target, convey.ShouldResemble, &testCopyStruct{
					Id:     100,
					Name:   "abc",
					Age:    23,
					Female: true,
				})
			})

			c.Convey("**struct to **struct", func(c convey.C) {
				var (
					target *testCopyStruct
					source *testCopyStruct = &testCopyStruct{
						Id:     100,
						Ignore: true,
						Name:   "abc",
						Age:    23,
						Female: true,
					}
				)

				err := copyFunc(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(target, convey.ShouldResemble, &testCopyStruct{
					Id:     100,
					Name:   "abc",
					Age:    23,
					Female: true,
				})
			})

			c.Convey("struct to *struct 2", func(c convey.C) {
				var (
					target testCopyStruct2
					source = testCopyStruct{
						Id:        100,
						Ignore:    true,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					}
				)

				c.Convey("with tag name", func(c convey.C) {
					err := copyFunc(&target, source, &CopyConfig{
						TargetTag: "gorm",
					})
					c.So(err, convey.ShouldBeNil)
					c.So(&target, convey.ShouldResemble, &testCopyStruct2{
						Id:        100,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					})
				})

				c.Convey("without tag name", func(c convey.C) {
					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(&target, convey.ShouldNotResemble, &testCopyStruct2{
						Id:        100,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					})
				})
			})

			c.Convey("*struct to *struct 2", func(c convey.C) {
				var (
					target testCopyStruct2
					source *testCopyStruct = &testCopyStruct{
						Id:        100,
						Ignore:    true,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					}
				)

				c.Convey("with tag name", func(c convey.C) {
					err := copyFunc(&target, source, &CopyConfig{
						TargetTag: "gorm",
					})
					c.So(err, convey.ShouldBeNil)
					c.So(&target, convey.ShouldResemble, &testCopyStruct2{
						Id:        100,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					})
				})

				c.Convey("without tag name", func(c convey.C) {
					err := copyFunc(&target, source)
					c.So(err, convey.ShouldBeNil)
					c.So(&target, convey.ShouldNotResemble, &testCopyStruct2{
						Id:        100,
						Name:      "abc",
						Age:       23,
						Female:    true,
						unexpored: new(int),
					})
				})
			})

			c.Convey("map[string]string to *struct", func(c convey.C) {
				var (
					target testCopyStruct
					source map[string]string = map[string]string{
						`id`:     "100",
						`Ignore`: "true",
						`Name`:   "abc",
						`Age`:    "23",
						`Female`: "true",
					}
				)

				err := copyFunc(&target, &source)
				c.So(err, convey.ShouldBeNil)
				c.So(&target, convey.ShouldResemble, &testCopyStruct{
					Id:     100,
					Name:   "abc",
					Age:    23,
					Female: true,
				})
			})

		})
	})
}

func TestDeepCopy(t *testing.T)    { testDeepCopyFunc(t, DeepCopy) }
func Test_deepCopyV1(t *testing.T) { testDeepCopyFunc(t, deepCopyV1) }
func Test_deepCopyV2(t *testing.T) { testDeepCopyFunc(t, deepCopyV2) }

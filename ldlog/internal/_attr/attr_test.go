/*
 * Copyright (C) distroy
 */

package _attr

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func testNewPtr[T any](v T) *T { return &v }

func TestAny(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		log := newTestLog(buf)
		getLogString := func() string {
			b := testGetLogString(buf, log)
			b = testRemoveLogPrefix(b)
			return b2s(b)
		}

		c.Convey("bool", func(c convey.C) {
			type Type = bool
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(true)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":true}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](false)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":false}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{false, true}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[false,true]}`)
				})
			})
		})

		c.Convey("int", func(c convey.C) {
			type Type = int
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("int8", func(c convey.C) {
			type Type = int8
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("int16", func(c convey.C) {
			type Type = int16
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("int32", func(c convey.C) {
			type Type = int32
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("int64", func(c convey.C) {
			type Type = int64
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})

		c.Convey("uint", func(c convey.C) {
			type Type = uint
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("uint8", func(c convey.C) {
			type Type = uint8
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			// c.Convey("slice", func(c convey.C) {
			// 	c.Convey("null", func(c convey.C) {
			// 		log.Debug("test", Any("key", ([]Type)(nil)))
			// 		s := getLogString()
			// 		c.So(s, convey.ShouldEqual, `{"key":null}`)
			// 	})
			// 	c.Convey("zero", func(c convey.C) {
			// 		log.Debug("test", Any("key", []Type{}))
			// 		s := getLogString()
			// 		c.So(s, convey.ShouldEqual, `{"key":[]}`)
			// 	})
			// 	c.Convey("valid", func(c convey.C) {
			// 		log.Debug("test", Any("key", []Type{0, 111}))
			// 		s := getLogString()
			// 		c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
			// 	})
			// })
		})
		c.Convey("uint16", func(c convey.C) {
			type Type = uint16
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("uint32", func(c convey.C) {
			type Type = uint32
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("uint64", func(c convey.C) {
			type Type = uint64
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("uintptr", func(c convey.C) {
			type Type = uintptr
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})

		c.Convey("float32", func(c convey.C) {
			type Type = float32
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})
		c.Convey("float64", func(c convey.C) {
			type Type = float64
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":123}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":0}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[0,111]}`)
				})
			})
		})

		c.Convey("complex64", func(c convey.C) {
			type Type = complex64
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(complex(123, 1))))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":"123+1i"}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"0+0i"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["0+0i","111+0i"]}`)
				})
			})
		})
		c.Convey("complex128", func(c convey.C) {
			type Type = complex128
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(complex(123, 1))))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":"123+1i"}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"0+0i"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 111}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["0+0i","111+0i"]}`)
				})
			})
		})

		c.Convey("duration", func(c convey.C) {
			type Type = time.Duration
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type(123*time.Millisecond)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":"123ms"}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](0)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"0s"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{0, 1023 * time.Millisecond}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["0s","1.023s"]}`)
				})
			})
		})
		c.Convey("time", func(c convey.C) {
			type Type = time.Time
			tz := time.FixedZone("Asia/Bejing", int(+(8*time.Hour)/time.Second))
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", time.Unix(1644479966, 0).In(tz)))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":"2022-02-10T15:59:26.000+0800"}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type](time.Time{})))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"0001-01-01T00:00:00.000Z"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{time.Unix(1644479966, 0).In(tz)}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["2022-02-10T15:59:26.000+0800"]}`)
				})
			})
		})

		c.Convey("string", func(c convey.C) {
			type Type = string
			c.Convey("value", func(c convey.C) {
				log.Debug("test", Any("key", Type("abc")))
				s := getLogString()
				c.So(s, convey.ShouldEqual, `{"key":"abc"}`)
			})
			c.Convey("pointer", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (*Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", testNewPtr[Type]("xyz")))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"xyz"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{"", "ijk"}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["","ijk"]}`)
				})
			})
		})
		c.Convey("byte string", func(c convey.C) {
			type Type = []byte
			c.Convey("value", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("nil", Any("key", Type(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":""}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", Type("abc")))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"abc"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{Type(""), Type("ijk")}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["","ijk"]}`)
				})
			})
		})
		c.Convey("stringer", func(c convey.C) {
			type Type = fmt.Stringer
			c.Convey("value", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", (Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", Type(string_t("abc"))))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"abc"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{string_t(""), string_t("ijk")}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":["","ijk"]}`)
				})
			})
		})

		c.Convey("error", func(c convey.C) {
			type Type = error
			c.Convey("null", func(c convey.C) {
				c.Convey("value", func(c convey.C) {
					log.Debug("nil", Any("key", Type(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", Type(fmt.Errorf("unknown error"))))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":"unknown error"}`)
				})
			})
			c.Convey("slice", func(c convey.C) {
				c.Convey("null", func(c convey.C) {
					log.Debug("test", Any("key", ([]Type)(nil)))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":null}`)
				})
				c.Convey("zero", func(c convey.C) {
					log.Debug("test", Any("key", []Type{}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[]}`)
				})
				c.Convey("valid", func(c convey.C) {
					log.Debug("test", Any("key", []Type{nil, fmt.Errorf("unknown error")}))
					s := getLogString()
					c.So(s, convey.ShouldEqual, `{"key":[null,"unknown error"]}`)
				})
			})
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldptr"
	"github.com/smartystreets/goconvey/convey"
)

func TestJsonCodec(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-json-codec"
		expiration := time.Duration(0)

		c.Convey("int64", func(c convey.C) {
			val := int64(100)

			set := New(rds, Json[any]()).Set(ctx, key, val, expiration)
			c.So(set.Err(), convey.ShouldBeNil)

			c.Convey("get-str", func(c convey.C) {
				get := rds.Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, "100")
			})
			c.Convey("get-i64", func(c convey.C) {
				get := New(rds, Json[int64]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, int64(100))
			})
			c.Convey("get-pi64", func(c convey.C) {
				get := New(rds, Json[*int64]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, ldptr.New[int64](100))
			})
			c.Convey("get-nil", func(c convey.C) {
				get := New(rds, Json[any]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, float64(100))
			})
		})
		c.Convey("object", func(c convey.C) {
			type Object struct {
				S    string  `json:"s"`
				PS   *string `json:"ps"`
				I64  int64   `json:"i64"`
				PI64 *int64  `json:"pi64"`
			}

			val := Object{
				S:    "abc",
				PS:   ldptr.New("xyz"),
				I64:  123,
				PI64: ldptr.New[int64](234),
			}

			set := New(rds, Json[any]()).Set(ctx, key, val, expiration)
			c.So(set.Err(), convey.ShouldBeNil)

			c.Convey("get-str", func(c convey.C) {
				get := rds.Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, `{"s":"abc","ps":"xyz","i64":123,"pi64":234}`)
			})
			c.Convey("get-obj", func(c convey.C) {
				get := New(rds, Json[Object]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, Object{
					S:    "abc",
					PS:   ldptr.New("xyz"),
					I64:  123,
					PI64: ldptr.New[int64](234),
				})
			})
			c.Convey("get-ptr", func(c convey.C) {
				get := New(rds, Json[*Object]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, &Object{
					S:    "abc",
					PS:   ldptr.New("xyz"),
					I64:  123,
					PI64: ldptr.New[int64](234),
				})
			})
			c.Convey("get-nil", func(c convey.C) {
				get := New(rds, Json[any]()).Get(ctx, key)
				c.So(get.Err(), convey.ShouldBeNil)
				c.So(get.Val(), convey.ShouldResemble, map[string]interface{}{
					"s":    "abc",
					"ps":   "xyz",
					"i64":  float64(123),
					"pi64": float64(234),
				})
			})
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldptr"
	"github.com/distroy/ldgo/v3/proto/ldtestproto"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

var (
	_ protoadapt.MessageV1 = (*ldtestproto.TestProto)(nil)
	_ protoadapt.MessageV2 = protoadapt.MessageV2Of(&ldtestproto.TestProto{})
)

func testMustProtoV1Marsha(m protoadapt.MessageV1) string {
	mm := protoadapt.MessageV2Of(m)
	b, _ := proto.Marshal(mm)
	return ldconv.BytesToStrUnsafe(b)
}

func testMustProtoV2Marsha(m protoadapt.MessageV2) string {
	b, _ := proto.Marshal(m)
	return ldconv.BytesToStrUnsafe(b)
}

func testClearUnexportField(v interface{}) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Interface {
		vv = vv.Elem()
	}

	if vv.Kind() != reflect.Ptr {
		return
	}
	vv = vv.Elem()

	if vv.Kind() != reflect.Struct {
		return
	}

	vt := vv.Type()
	for i, n := 0, vv.NumField(); i < n; i++ {
		ff := vt.Field(i)
		if ff.IsExported() {
			continue
		}

		f := vv.Field(i)

		a := unsafe.Pointer(f.UnsafeAddr())
		o := reflect.NewAt(f.Type(), a).Elem()
		o.Set(reflect.Zero(f.Type()))
	}
}

func TestProtoV1Codec(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-proto-v1-codec"
		expiration := time.Duration(0)
		val := &ldtestproto.TestProto{
			Str: ldptr.New("abc"),
			I64: ldptr.New[int64](234),
			F64: ldptr.New(100.234),
		}

		set := New(rds, ProtoV1[protoadapt.MessageV1]()).Set(ctx, key, val, expiration)
		c.So(set.Err(), convey.ShouldBeNil)

		c.Convey("get-str", func(c convey.C) {
			get := rds.Get(ctx, key)
			c.So(get.Err(), convey.ShouldBeNil)
			c.So(get.Val(), convey.ShouldResemble, testMustProtoV1Marsha(val))
		})
		c.Convey("get-obj", func(c convey.C) {
			get := New(rds, ProtoV1[*ldtestproto.TestProto]()).Get(ctx, key)
			c.So(get.Err(), convey.ShouldBeNil)

			testClearUnexportField(get.Val())
			c.So(get.Val(), convey.ShouldResemble, &ldtestproto.TestProto{
				Str: ldptr.New("abc"),
				I64: ldptr.New[int64](234),
				F64: ldptr.New(100.234),
			})
		})
	})
}

func TestProtoV2Codec(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-proto-v2-codec"
		expiration := time.Duration(0)
		val := &ldtestproto.TestProto{
			Str: ldptr.New("abc"),
			I64: ldptr.New[int64](234),
			F64: ldptr.New(100.234),
		}

		set := New(rds, ProtoV2[protoadapt.MessageV2]()).Set(ctx, key, val, expiration)
		c.So(set.Err(), convey.ShouldBeNil)

		c.Convey("get-str", func(c convey.C) {
			get := rds.Get(ctx, key)
			c.So(get.Err(), convey.ShouldBeNil)
			c.So(get.Val(), convey.ShouldResemble, testMustProtoV2Marsha(val))
		})
		c.Convey("get-obj", func(c convey.C) {
			get := New(rds, ProtoV2[*ldtestproto.TestProto]()).Get(ctx, key)
			c.So(get.Err(), convey.ShouldBeNil)

			testClearUnexportField(get.Val())
			c.So(get.Val(), convey.ShouldResemble, &ldtestproto.TestProto{
				Str: ldptr.New("abc"),
				I64: ldptr.New[int64](234),
				F64: ldptr.New(100.234),
			})
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/distroy/ldgo/v3/ldredis"
	"github.com/distroy/ldgo/v3/ldredis/internal"
)

type Codec[T any] interface {
	Marshal(v T) ([]byte, error)
	Unmarshal(b []byte) (T, error)
}

type base[T any] struct {
	*client

	codec Codec[T]
}

func (c *base[T]) zero() T {
	var zero T
	return zero
}

// ***** redis marshal begin *****

func (c *base[T]) marshal(value T) (interface{}, error) {
	// if value == nil {
	// 	return nil, nil
	// }
	// v, ok := value.(T)
	// if !ok {
	// 	err := fmt.Errorf("redis codec cannot marshal the type `%T`", value)
	// 	return errorMarshaler{Err: err}, err
	// }
	bytes, err := c.codec.Marshal(value)
	if err != nil {
		return errorMarshaler{Err: err}, err
	}
	return ldconv.BytesToStrUnsafe(bytes), nil
}

func (c *base[T]) mustMarshal(value T) interface{} {
	res, _ := c.marshal(value)
	return res
}

func (c *base[T]) _marshalPairs(pairs []interface{}) []interface{} {
	for i, l := 1, len(pairs); i < l; i += 2 {
		v, ok := pairs[i].(T)
		if !ok {
			err := fmt.Errorf("redis codec cannot marshal the type `%T`", pairs[i])
			pairs[i] = errorMarshaler{Err: err}
			break
		}

		val, err := c.marshal(v)
		pairs[i] = val
		if err != nil {
			break
		}
	}
	return pairs
}

func (c *base[T]) marshalPairs(pairs []Pair[T]) []interface{} {
	res := make([]interface{}, 0, len(pairs)*2)
	for _, v := range pairs {
		res = append(res, v.First)
		val, err := c.marshal(v.Second)
		res = append(res, val)
		if err != nil {
			break
		}
	}
	return res
}

func (c *base[T]) marshalSlice(s []T) []interface{} {
	r := make([]interface{}, 0, len(s))
	for _, v := range s {
		val, err := c.marshal(v)
		r = append(r, val)
		if err != nil {
			break
		}
	}
	return r
}

func (c *base[T]) marshalMap(m map[string]T) []interface{} {
	r := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		val, err := c.marshal(v)
		r = append(r, k, val)
		if err != nil {
			break
		}
	}
	return r
}

func (c *base[T]) marshalZMember(v ZMember[T]) ldredis.ZMember {
	return ldredis.ZMember{
		Score:  v.Score,
		Member: c.mustMarshal(v.Member),
	}
}
func (c *base[T]) marshalZMembers(s []ZMember[T]) []ldredis.ZMember {
	r := make([]ldredis.ZMember, 0, len(s))
	for _, v := range s {
		val, err := c.marshal(v.Member)
		r = append(r, ldredis.ZMember{
			Score:  v.Score,
			Member: val,
		})
		if err != nil {
			break
		}
	}
	return r
}

// ***** redis marshal end *****

// ***** redis unmarshal begin *****

func (c *base[T]) unmarshal(cc context.Context, bytes []byte) (T, error) {
	if len(bytes) == 0 {
		return c.zero(), nil
	}
	v, err := c.codec.Unmarshal(bytes)
	if err != nil {
		logger := ldctx.GetLogger(cc)
		logger.Error("redis codec unmarshal fail", ldlog.Error(err),
			getCallerField(c.client))
		return c.zero(), lderr.ErrCacheUnmarshal
	}
	return v, nil
}

func (c *base[T]) unmarshalInterface(cc context.Context, i interface{}) (T, error) {
	if i == nil {
		return c.zero(), nil
	}

	var bytes []byte
	switch tmp := i.(type) {
	case []byte:
		bytes = tmp
	case string:
		bytes = ldconv.StrToBytesUnsafe(tmp)
	default:
		logger := ldctx.GetLogger(cc)
		logger.Error("redis codec cannot unmarshal the type",
			ldlog.Stringer("type", reflect.TypeOf(i)), getCallerField(c.client))
		return c.zero(), lderr.ErrCacheUnmarshal
	}
	return c.unmarshal(cc, bytes)
}

// ***** redis unmarshal end *****

func getCallerField(rds *ldredis.Redis) ldlog.Attr {
	caller := true
	if rds != nil {
		caller = getOptions(rds).Caller
	}
	return internal.GetCallerField(caller)
}

//go:linkname getOptions github.com/distroy/ldgo/v3/ldredis.getOptions
func getOptions(*ldredis.Redis) *internal.Options

//go:linkname getOptionsPointer github.com/distroy/ldgo/v3/ldredis.getOptionsPointer
func getOptionsPointer(c *ldredis.Redis) unsafe.Pointer

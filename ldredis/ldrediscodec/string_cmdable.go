/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldredis"
)

type StringCmdable[T any] interface {
	// ldredis.StringCmdable

	Append(cc context.Context, key, value string) *ldredis.IntCmd                               // same as ldredis.Cmdable
	Decr(cc context.Context, key string) *ldredis.IntCmd                                        // same as ldredis.Cmdable
	DecrBy(cc context.Context, key string, decrement int64) *ldredis.IntCmd                     // same as ldredis.Cmdable
	Get(cc context.Context, key string) *TypeCmd[T]                                             //
	GetRange(cc context.Context, key string, start, end int64) *ldredis.StringCmd               // same as ldredis.Cmdable
	GetSet(cc context.Context, key string, value T) *TypeCmd[T]                                 //
	GetEx(cc context.Context, key string, expiration time.Duration) *ldredis.StringCmd          // same as ldredis.Cmdable
	GetDel(cc context.Context, key string) *ldredis.StringCmd                                   // same as ldredis.Cmdable
	Incr(cc context.Context, key string) *ldredis.IntCmd                                        // same as ldredis.Cmdable
	IncrBy(cc context.Context, key string, value int64) *ldredis.IntCmd                         // same as ldredis.Cmdable
	IncrByFloat(cc context.Context, key string, value float64) *ldredis.FloatCmd                // same as ldredis.Cmdable
	MGet(cc context.Context, keys ...string) *SliceCmd[T]                                       //
	MSet(cc context.Context, key string, value T) *ldredis.StatusCmd                            //
	MSetPairs(cc context.Context, pairs ...Pair[T]) *ldredis.StatusCmd                          //
	MSetMap(cc context.Context, KeyVals map[string]T) *ldredis.StatusCmd                        //
	MSetNX(cc context.Context, key string, value T) *ldredis.BoolCmd                            //
	MSetNXPairs(cc context.Context, pairs ...Pair[T]) *ldredis.BoolCmd                          //
	MSetNXMap(cc context.Context, KeyVals map[string]T) *ldredis.BoolCmd                        //
	Set(cc context.Context, key string, value T, expiration time.Duration) *ldredis.StatusCmd   //
	SetArgs(cc context.Context, key string, value T, a ldredis.SetArgs) *ldredis.StatusCmd      //
	SetEx(cc context.Context, key string, value T, expiration time.Duration) *ldredis.StatusCmd //
	SetNX(cc context.Context, key string, value T, expiration time.Duration) *ldredis.BoolCmd   //
	SetXX(cc context.Context, key string, value T, expiration time.Duration) *ldredis.BoolCmd   //
	SetRange(cc context.Context, key string, offset int64, value string) *ldredis.IntCmd        // same as ldredis.Cmdable
	StrLen(cc context.Context, key string) *ldredis.IntCmd                                      // same as ldredis.Cmdable
}

func (c *Redis[T]) MGet(cc context.Context, keys ...string) *SliceCmd[T] {
	return newSliceCmd(cc, c.base, c.client.MGet(cc, keys...))
}
func (c *Redis[T]) MSet(cc context.Context, key string, value T) *ldredis.StatusCmd {
	return c.client.MSet(cc, key, c.mustMarshal(value))
}
func (c *Redis[T]) MSetPairs(cc context.Context, pairs ...Pair[T]) *ldredis.StatusCmd {
	return c.client.MSet(cc, c.marshalPairs(pairs)...)
}
func (c *Redis[T]) MSetMap(cc context.Context, keyVals map[string]T) *ldredis.StatusCmd {
	return c.client.MSet(cc, c.marshalMap(keyVals)...)
}
func (c *Redis[T]) MSetNX(cc context.Context, key string, value T) *ldredis.BoolCmd {
	return c.client.MSetNX(cc, key, c.mustMarshal(value))
}
func (c *Redis[T]) MSetNXPairs(cc context.Context, pairs ...Pair[T]) *ldredis.BoolCmd {
	return c.client.MSetNX(cc, c.marshalPairs(pairs)...)
}
func (c *Redis[T]) MSetNXMap(cc context.Context, keyVals map[string]T) *ldredis.BoolCmd {
	return c.client.MSetNX(cc, c.marshalMap(keyVals)...)
}

func (c *Redis[T]) SetArgs(cc context.Context, key string, value T, a ldredis.SetArgs) *ldredis.StatusCmd {
	return c.client.SetArgs(cc, key, c.mustMarshal(value), a)
}

func (c *Redis[T]) Get(cc context.Context, key string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.Get(cc, key))
}
func (c *Redis[T]) GetSet(cc context.Context, key string, value T) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.GetSet(cc, key, c.mustMarshal(value)))
}
func (c *Redis[T]) Set(cc context.Context, key string, value T, expiration time.Duration) *ldredis.StatusCmd {
	return c.client.Set(cc, key, c.mustMarshal(value), expiration)
}
func (c *Redis[T]) SetEx(cc context.Context, key string, value T, expiration time.Duration) *ldredis.StatusCmd {
	return c.client.SetEx(cc, key, c.mustMarshal(value), expiration)
}
func (c *Redis[T]) SetNX(cc context.Context, key string, value T, expiration time.Duration) *ldredis.BoolCmd {
	return c.client.SetNX(cc, key, c.mustMarshal(value), expiration)
}
func (c *Redis[T]) SetXX(cc context.Context, key string, value T, expiration time.Duration) *ldredis.BoolCmd {
	return c.client.SetXX(cc, key, c.mustMarshal(value), expiration)
}

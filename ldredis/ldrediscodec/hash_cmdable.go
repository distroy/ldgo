/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"

	"github.com/distroy/ldgo/v3/ldredis"
)

type HashCmdable[T comparable] interface {
	HDel(cc context.Context, key string, fields ...string) *ldredis.IntCmd                           // same as ldredis.Cmdable
	HExists(cc context.Context, key, field string) *ldredis.BoolCmd                                  // same as ldredis.Cmdable
	HGet(cc context.Context, key, field string) *TypeCmd[T]                                          //
	HGetAll(cc context.Context, key string) *MapStringTypeCmd[T]                                     //
	HIncrBy(cc context.Context, key, field string, incr int64) *ldredis.IntCmd                       // same as ldredis.Cmdable
	HIncrByFloat(cc context.Context, key, field string, incr float64) *ldredis.FloatCmd              // same as ldredis.Cmdable
	HKeys(cc context.Context, key string) *ldredis.StringSliceCmd                                    // same as ldredis.Cmdable
	HLen(cc context.Context, key string) *ldredis.IntCmd                                             // same as ldredis.Cmdable
	HMGet(cc context.Context, key string, fields ...string) *SliceCmd[T]                             //
	HSet(cc context.Context, key string, field string, value T) *ldredis.IntCmd                      //
	HSetPairs(cc context.Context, key string, fields ...Pair[T]) *ldredis.IntCmd                     //
	HSetMap(cc context.Context, key string, fields map[string]T) *ldredis.IntCmd                     //
	HMSet(cc context.Context, key string, field string, value T) *ldredis.BoolCmd                    //
	HMSetPairs(cc context.Context, key string, fields ...Pair[T]) *ldredis.BoolCmd                   //
	HMSetMap(cc context.Context, key string, fields map[string]T) *ldredis.BoolCmd                   //
	HSetNX(cc context.Context, key, field string, value T) *ldredis.BoolCmd                          //
	HScan(cc context.Context, key string, cursor uint64, match string, count int64) *ldredis.ScanCmd // same as ldredis.Cmdable
	HVals(cc context.Context, key string) *TypeSliceCmd[T]                                           //
	HRandField(cc context.Context, key string, count int) *ldredis.StringSliceCmd                    // same as ldredis.Cmdable
	HRandFieldWithValues(cc context.Context, key string, count int) *ldredis.KeyValueSliceCmd        // same as ldredis.Cmdable
}

func (c *Redis[T]) HGet(cc context.Context, key, field string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.HGet(cc, key, field))
}
func (c *Redis[T]) HGetAll(cc context.Context, key string) *MapStringTypeCmd[T] {
	return newMapStringTypeCmd[T](cc, c.base, c.client.HGetAll(cc, key))
}
func (c *Redis[T]) HMGet(cc context.Context, key string, fields ...string) *SliceCmd[T] {
	return newSliceCmd(cc, c.base, c.client.HMGet(cc, key, fields...))
}
func (c *Redis[T]) HMSet(cc context.Context, key, field string, value T) *ldredis.BoolCmd {
	return c.client.HMSet(cc, key, field, c.mustMarshal(value))
}
func (c *Redis[T]) HMSetPairs(cc context.Context, key string, pairs ...Pair[T]) *ldredis.BoolCmd {
	return c.client.HMSet(cc, key, c.marshalPairs(pairs)...)
}
func (c *Redis[T]) HMSetMap(cc context.Context, key string, fields map[string]T) *ldredis.BoolCmd {
	return c.client.HMSet(cc, key, c.marshalMap(fields)...)
}
func (c *Redis[T]) HSet(cc context.Context, key, field string, value T) *ldredis.IntCmd {
	return c.client.HSet(cc, key, field, c.mustMarshal(value))
}
func (c *Redis[T]) HSetPairs(cc context.Context, key string, pairs ...Pair[T]) *ldredis.IntCmd {
	return c.client.HSet(cc, key, c.marshalPairs(pairs)...)
}
func (c *Redis[T]) HSetMap(cc context.Context, key string, fields map[string]T) *ldredis.IntCmd {
	return c.client.HSet(cc, key, c.marshalMap(fields)...)
}
func (c *Redis[T]) HSetNX(cc context.Context, key, field string, value T) *ldredis.BoolCmd {
	return c.client.HSetNX(cc, key, field, c.mustMarshal(value))
}
func (c *Redis[T]) HVals(cc context.Context, key string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.HVals(cc, key))
}

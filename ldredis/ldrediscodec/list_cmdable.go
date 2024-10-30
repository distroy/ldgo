/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldredis"
)

type ListCmdable[T any] interface {
	BLPop(cc context.Context, timeout time.Duration, keys ...string) *TypeSliceCmd[T]                                      //
	BLMPop(cc context.Context, timeout time.Duration, direction string, count int64, keys ...string) *ldredis.KeyValuesCmd //
	BRPop(cc context.Context, timeout time.Duration, keys ...string) *TypeSliceCmd[T]                                      //
	BRPopLPush(cc context.Context, source, destination string, timeout time.Duration) *TypeCmd[T]                          //
	LIndex(cc context.Context, key string, index int64) *TypeCmd[T]                                                        //
	LInsert(cc context.Context, key, op string, pivot, value T) *ldredis.IntCmd                                            //
	LInsertBefore(cc context.Context, key string, pivot, value T) *ldredis.IntCmd                                          //
	LInsertAfter(cc context.Context, key string, pivot, value T) *ldredis.IntCmd                                           //
	LLen(cc context.Context, key string) *ldredis.IntCmd                                                                   // same as ldredis.Cmdable
	LMPop(cc context.Context, direction string, count int64, keys ...string) *ldredis.KeyValuesCmd                         //
	LPop(cc context.Context, key string) *TypeCmd[T]                                                                       //
	LPopCount(cc context.Context, key string, count int) *ldredis.StringSliceCmd                                           //
	LPos(cc context.Context, key string, value string, args ldredis.LPosArgs) *ldredis.IntCmd                              //
	LPosCount(cc context.Context, key string, value string, count int64, args ldredis.LPosArgs) *ldredis.IntSliceCmd       //
	LPush(cc context.Context, key string, values ...T) *ldredis.IntCmd                                                     //
	LPushX(cc context.Context, key string, values ...T) *ldredis.IntCmd                                                    //
	LRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T]                                             //
	LRem(cc context.Context, key string, count int64, value T) *ldredis.IntCmd                                             //
	LSet(cc context.Context, key string, index int64, value T) *ldredis.StatusCmd                                          //
	LTrim(cc context.Context, key string, start, stop int64) *ldredis.StatusCmd                                            // same as ldredis.Cmdable
	RPop(cc context.Context, key string) *TypeCmd[T]                                                                       //
	RPopCount(cc context.Context, key string, count int) *ldredis.StringSliceCmd                                           //
	RPopLPush(cc context.Context, source, destination string) *TypeCmd[T]                                                  //
	RPush(cc context.Context, key string, values ...T) *ldredis.IntCmd                                                     //
	RPushX(cc context.Context, key string, values ...T) *ldredis.IntCmd                                                    //
	LMove(cc context.Context, source, destination, srcpos, destpos string) *ldredis.StringCmd                              //
	BLMove(cc context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *ldredis.StringCmd      //
}

func (c *Redis[T]) BLPop(cc context.Context, timeout time.Duration, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.BLPop(cc, timeout, keys...))
}
func (c *Redis[T]) BRPop(cc context.Context, timeout time.Duration, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.BRPop(cc, timeout, keys...))
}
func (c *Redis[T]) BRPopLPush(cc context.Context, source, destination string, timeout time.Duration) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.BRPopLPush(cc, source, destination, timeout))
}
func (c *Redis[T]) LIndex(cc context.Context, key string, index int64) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.LIndex(cc, key, index))
}
func (c *Redis[T]) LInsert(cc context.Context, key, op string, pivot, value T) *ldredis.IntCmd {
	return c.client.LInsert(cc, key, op, c.mustMarshal(pivot), c.mustMarshal(value))
}
func (c *Redis[T]) LInsertBefore(cc context.Context, key string, pivot, value T) *ldredis.IntCmd {
	return c.client.LInsertBefore(cc, key, c.mustMarshal(pivot), c.mustMarshal(value))
}
func (c *Redis[T]) LInsertAfter(cc context.Context, key string, pivot, value T) *ldredis.IntCmd {
	return c.client.LInsertAfter(cc, key, c.mustMarshal(pivot), c.mustMarshal(value))
}
func (c *Redis[T]) LPop(cc context.Context, key string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.LPop(cc, key))
}
func (c *Redis[T]) LPush(cc context.Context, key string, values ...T) *ldredis.IntCmd {
	return c.client.LPush(cc, key, c.marshalSlice(values)...)
}
func (c *Redis[T]) LPushX(cc context.Context, key string, values ...T) *ldredis.IntCmd {
	return c.client.LPushX(cc, key, c.marshalSlice(values)...)
}
func (c *Redis[T]) LRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.LRange(cc, key, start, stop))
}
func (c *Redis[T]) LRem(cc context.Context, key string, count int64, value T) *ldredis.IntCmd {
	return c.client.LRem(cc, key, count, c.mustMarshal(value))
}
func (c *Redis[T]) LSet(cc context.Context, key string, index int64, value T) *ldredis.StatusCmd {
	return c.client.LSet(cc, key, index, c.mustMarshal(value))
}
func (c *Redis[T]) RPop(cc context.Context, key string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.RPop(cc, key))
}
func (c *Redis[T]) RPopLPush(cc context.Context, source, destination string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.RPopLPush(cc, source, destination))
}
func (c *Redis[T]) RPush(cc context.Context, key string, values ...T) *ldredis.IntCmd {
	return c.client.RPush(cc, key, c.marshalSlice(values)...)
}
func (c *Redis[T]) RPushX(cc context.Context, key string, values ...T) *ldredis.IntCmd {
	return c.client.RPushX(cc, key, c.marshalSlice(values)...)
}

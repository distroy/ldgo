/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"

	"github.com/distroy/ldgo/v3/ldredis"
)

type SetCmdable[T comparable] interface {
	SAdd(cc context.Context, key string, members ...T) *ldredis.IntCmd                               //
	SCard(cc context.Context, key string) *ldredis.IntCmd                                            // same as ldredis.Cmdable
	SDiff(cc context.Context, keys ...string) *TypeSliceCmd[T]                                       //
	SDiffStore(cc context.Context, destination string, keys ...string) *ldredis.IntCmd               // same as ldredis.Cmdable
	SInter(cc context.Context, keys ...string) *TypeSliceCmd[T]                                      //
	SInterCard(cc context.Context, limit int64, keys ...string) *ldredis.IntCmd                      // same as ldredis.Cmdable
	SInterStore(cc context.Context, destination string, keys ...string) *ldredis.IntCmd              // same as ldredis.Cmdable
	SIsMember(cc context.Context, key string, member T) *ldredis.BoolCmd                             //
	SMIsMember(cc context.Context, key string, members ...T) *ldredis.BoolSliceCmd                   //
	SMembers(cc context.Context, key string) *TypeSliceCmd[T]                                        //
	SMembersMap(cc context.Context, key string) *TypeSetCmd[T]                                       //
	SMove(cc context.Context, source, destination string, member T) *ldredis.BoolCmd                 //
	SPop(cc context.Context, key string) *TypeCmd[T]                                                 //
	SPopN(cc context.Context, key string, count int64) *TypeSliceCmd[T]                              //
	SRandMember(cc context.Context, key string) *TypeCmd[T]                                          //
	SRandMemberN(cc context.Context, key string, count int64) *TypeSliceCmd[T]                       //
	SRem(cc context.Context, key string, members ...T) *ldredis.IntCmd                               //
	SScan(cc context.Context, key string, cursor uint64, match string, count int64) *ldredis.ScanCmd // same as ldredis.Cmdable
	SUnion(cc context.Context, keys ...string) *TypeSliceCmd[T]                                      //
	SUnionStore(cc context.Context, destination string, keys ...string) *ldredis.IntCmd              // same as ldredis.Cmdable
}

func (c *Redis[T]) SAdd(cc context.Context, key string, members ...T) *ldredis.IntCmd {
	return c.client.SAdd(cc, key, c.marshalSlice(members)...)
}
func (c *Redis[T]) SDiff(cc context.Context, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SDiff(cc, keys...))
}

func (c *Redis[T]) SInter(cc context.Context, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SInter(cc, keys...))
}

func (c *Redis[T]) SIsMember(cc context.Context, key string, member T) *ldredis.BoolCmd {
	return c.client.SIsMember(cc, key, c.mustMarshal(member))
}
func (c *Redis[T]) SMIsMember(cc context.Context, key string, members ...T) *ldredis.BoolSliceCmd {
	return c.client.SMIsMember(cc, key, c.marshalSlice(members)...)
}
func (c *Redis[T]) SMembers(cc context.Context, key string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SMembers(cc, key))
}
func (c *Redis[T]) SMembersMap(cc context.Context, key string) *TypeSetCmd[T] {
	return newTypeSetCmd[T](cc, c.base, c.client.SMembersMap(cc, key))
}
func (c *Redis[T]) SMove(cc context.Context, src, dest string, member T) *ldredis.BoolCmd {
	return c.client.SMove(cc, src, dest, c.mustMarshal(member))
}
func (c *Redis[T]) SPop(cc context.Context, key string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.SPop(cc, key))
}
func (c *Redis[T]) SPopN(cc context.Context, key string, count int64) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SPopN(cc, key, count))
}
func (c *Redis[T]) SRandMember(cc context.Context, key string) *TypeCmd[T] {
	return newTypeCmd[T](cc, c.base, c.client.SRandMember(cc, key))
}
func (c *Redis[T]) SRandMemberN(cc context.Context, key string, count int64) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SRandMemberN(cc, key, count))
}
func (c *Redis[T]) SRem(cc context.Context, key string, members ...T) *ldredis.IntCmd {
	return c.client.SRem(cc, key, c.marshalSlice(members)...)
}
func (c *Redis[T]) SUnion(cc context.Context, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.SUnion(cc, keys...))
}

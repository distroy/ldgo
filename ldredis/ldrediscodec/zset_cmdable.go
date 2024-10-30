/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldredis"
	"github.com/redis/go-redis/v9"
)

type SortedSetCmdable[T any] interface {
	BZPopMax(cc context.Context, timeout time.Duration, keys ...string) *ZMemberWithKeyCmd[T]                               //
	BZPopMin(cc context.Context, timeout time.Duration, keys ...string) *ZMemberWithKeyCmd[T]                               //
	BZMPop(cc context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZMemberSliceWithKeyCmd[T] //
	ZAdd(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd                                             //
	ZAddLT(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd                                           //
	ZAddGT(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd                                           //
	ZAddNX(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd                                           //
	ZAddXX(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd                                           //
	ZAddArgs(cc context.Context, key string, args ZAddArgs[T]) *ldredis.IntCmd                                              //
	ZAddArgsIncr(cc context.Context, key string, args ZAddArgs[T]) *ldredis.FloatCmd                                        //
	ZCard(cc context.Context, key string) *ldredis.IntCmd                                                                   // same as ldredis.Cmdable
	ZCount(cc context.Context, key, min, max string) *ldredis.IntCmd                                                        // same as ldredis.Cmdable
	ZLexCount(cc context.Context, key, min, max string) *ldredis.IntCmd                                                     // same as ldredis.Cmdable
	ZIncrBy(cc context.Context, key string, increment float64, member T) *ldredis.FloatCmd                                  //
	ZInter(cc context.Context, store *ldredis.ZStore) *TypeSliceCmd[T]                                                      //
	ZInterWithScores(cc context.Context, store *ldredis.ZStore) *ZMemberSliceCmd[T]                                         //
	ZInterCard(cc context.Context, limit int64, keys ...string) *ldredis.IntCmd                                             // same as ldredis.Cmdable
	ZInterStore(cc context.Context, destination string, store *ldredis.ZStore) *ldredis.IntCmd                              // same as ldredis.Cmdable
	ZMPop(cc context.Context, order string, count int64, keys ...string) *ZMemberSliceWithKeyCmd[T]                         //
	ZMScore(cc context.Context, key string, members ...T) *ldredis.FloatSliceCmd                                            //
	ZPopMax(cc context.Context, key string, count ...int64) *ZMemberSliceCmd[T]                                             //
	ZPopMin(cc context.Context, key string, count ...int64) *ZMemberSliceCmd[T]                                             //
	ZRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T]                                              //
	ZRangeWithScores(cc context.Context, key string, start, stop int64) *ZMemberSliceCmd[T]                                 //
	ZRangeByScore(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T]                                   //
	ZRangeByLex(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T]                                     //
	ZRangeByScoreWithScores(cc context.Context, key string, opt *ldredis.ZRangeBy) *ZMemberSliceCmd[T]                      //
	ZRangeArgs(cc context.Context, z ldredis.ZRangeArgs) *TypeSliceCmd[T]                                                   //
	ZRangeArgsWithScores(cc context.Context, z ldredis.ZRangeArgs) *ZMemberSliceCmd[T]                                      //
	ZRangeStore(cc context.Context, dstKey string, z ldredis.ZRangeArgs) *ldredis.IntCmd                                    // same as ldredis.Cmdable
	ZRank(cc context.Context, key string, member T) *ldredis.IntCmd                                                         //
	ZRankWithScore(cc context.Context, key, member string) *ldredis.RankWithScoreCmd                                        // same as ldredis.Cmdable
	ZRem(cc context.Context, key string, members ...T) *ldredis.IntCmd                                                      //
	ZRemRangeByRank(cc context.Context, key string, start, stop int64) *ldredis.IntCmd                                      // same as ldredis.Cmdable
	ZRemRangeByScore(cc context.Context, key, min, max string) *ldredis.IntCmd                                              // same as ldredis.Cmdable
	ZRemRangeByLex(cc context.Context, key, min, max string) *ldredis.IntCmd                                                // same as ldredis.Cmdable
	ZRevRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T]                                           //
	ZRevRangeWithScores(cc context.Context, key string, start, stop int64) *ZMemberSliceCmd[T]                              //
	ZRevRangeByScore(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T]                                //
	ZRevRangeByLex(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T]                                  //
	ZRevRangeByScoreWithScores(cc context.Context, key string, opt *ldredis.ZRangeBy) *ZMemberSliceCmd[T]                   //
	ZRevRank(cc context.Context, key string, member T) *ldredis.IntCmd                                                      //
	ZRevRankWithScore(cc context.Context, key string, member T) *ldredis.RankWithScoreCmd                                   //
	ZScore(cc context.Context, key string, member T) *ldredis.FloatCmd                                                      //
	ZUnionStore(cc context.Context, dest string, store *ldredis.ZStore) *ldredis.IntCmd                                     // same as ldredis.Cmdable
	ZRandMember(cc context.Context, key string, count int) *TypeSliceCmd[T]                                                 //
	ZRandMemberWithScores(cc context.Context, key string, count int) *ZMemberSliceCmd[T]                                    //
	ZUnion(cc context.Context, store ldredis.ZStore) *TypeSliceCmd[T]                                                       //
	ZUnionWithScores(cc context.Context, store ldredis.ZStore) *ZMemberSliceCmd[T]                                          //
	ZDiff(cc context.Context, keys ...string) *TypeSliceCmd[T]                                                              //
	ZDiffWithScores(cc context.Context, keys ...string) *ZMemberSliceCmd[T]                                                 //
	ZDiffStore(cc context.Context, destination string, keys ...string) *ldredis.IntCmd                                      // same as ldredis.Cmdable
	ZScan(cc context.Context, key string, cursor uint64, match string, count int64) *ldredis.ScanCmd                        // same as ldredis.Cmdable
}

func (c *Redis[T]) BZPopMax(cc context.Context, timeout time.Duration, keys ...string) *ZMemberWithKeyCmd[T] {
	return newZMemberWithKeyCmd[T](cc, c.base, c.client.BZPopMax(cc, timeout, keys...))
}
func (c *Redis[T]) BZPopMin(cc context.Context, timeout time.Duration, keys ...string) *ZMemberWithKeyCmd[T] {
	return newZMemberWithKeyCmd[T](cc, c.base, c.client.BZPopMin(cc, timeout, keys...))
}

func (c *Redis[T]) BZMPop(cc context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZMemberSliceWithKeyCmd[T] {
	return newZMemberSliceWithKeyCmd[T](cc, c.base, c.client.BZMPop(cc, timeout, order, count, keys...))
}

func (c *Redis[T]) ZAdd(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd {
	return c.client.ZAdd(cc, key, c.marshalZMembers(members)...)
}
func (c *Redis[T]) ZAddLT(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd {
	return c.client.ZAddLT(cc, key, c.marshalZMembers(members)...)
}
func (c *Redis[T]) ZAddGT(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd {
	return c.client.ZAddGT(cc, key, c.marshalZMembers(members)...)
}
func (c *Redis[T]) ZAddNX(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd {
	return c.client.ZAddNX(cc, key, c.marshalZMembers(members)...)
}
func (c *Redis[T]) ZAddXX(cc context.Context, key string, members ...ZMember[T]) *ldredis.IntCmd {
	return c.client.ZAddXX(cc, key, c.marshalZMembers(members)...)
}

func (c *Redis[T]) ZAddArgs(cc context.Context, key string, args ZAddArgs[T]) *ldredis.IntCmd {
	return c.client.ZAddArgs(cc, key, redis.ZAddArgs{
		NX:      args.NX,
		XX:      args.XX,
		LT:      args.LT,
		GT:      args.GT,
		Ch:      args.Ch,
		Members: c.marshalZMembers(args.Members),
	})
}
func (c *Redis[T]) ZAddArgsIncr(cc context.Context, key string, args ZAddArgs[T]) *ldredis.FloatCmd {
	return c.client.ZAddArgsIncr(cc, key, redis.ZAddArgs{
		NX:      args.NX,
		XX:      args.XX,
		LT:      args.LT,
		GT:      args.GT,
		Ch:      args.Ch,
		Members: c.marshalZMembers(args.Members),
	})
}

func (c *Redis[T]) ZPopMax(cc context.Context, key string, count ...int64) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZPopMax(cc, key, count...))
}
func (c *Redis[T]) ZIncrBy(cc context.Context, key string, increment float64, member T) *ldredis.FloatCmd {
	// return c.client.ZIncrBy(cc, key, increment, c.mustMarshal(member))
	cmd := redis.NewFloatCmd(cc, "zincrby", key, increment, c.mustMarshal(member))
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZInter(cc context.Context, store *ldredis.ZStore) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZInter(cc, store))
}
func (c *Redis[T]) ZInterWithScores(cc context.Context, store *ldredis.ZStore) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZInterWithScores(cc, store))
}
func (c *Redis[T]) ZMPop(cc context.Context, order string, count int64, keys ...string) *ZMemberSliceWithKeyCmd[T] {
	return newZMemberSliceWithKeyCmd[T](cc, c.base, c.client.ZMPop(cc, order, count, keys...))
}
func (c *Redis[T]) ZMScore(cc context.Context, key string, members ...T) *ldredis.FloatSliceCmd {
	// return c.client.ZMScore(cc, key, c.marshalSlice(members)...)
	args := []interface{}{"zmscore", key}
	args = append(args, c.marshalSlice(members)...)
	cmd := redis.NewFloatSliceCmd(cc, args...)
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZPopMin(cc context.Context, key string, count ...int64) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZPopMin(cc, key, count...))
}
func (c *Redis[T]) ZRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRange(cc, key, start, stop))
}
func (c *Redis[T]) ZRangeWithScores(cc context.Context, key string, start, stop int64) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRangeWithScores(cc, key, start, stop))
}
func (c *Redis[T]) ZRangeByScore(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRangeByScore(cc, key, opt))
}
func (c *Redis[T]) ZRangeByLex(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRangeByLex(cc, key, opt))
}
func (c *Redis[T]) ZRangeByScoreWithScores(cc context.Context, key string, opt *ldredis.ZRangeBy) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRangeByScoreWithScores(cc, key, opt))
}
func (c *Redis[T]) ZRangeArgs(cc context.Context, z ldredis.ZRangeArgs) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRangeArgs(cc, z))
}
func (c *Redis[T]) ZRangeArgsWithScores(cc context.Context, z ldredis.ZRangeArgs) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRangeArgsWithScores(cc, z))
}
func (c *Redis[T]) ZRank(cc context.Context, key string, member T) *ldredis.IntCmd {
	cmd := redis.NewIntCmd(cc, "zrank", key, c.mustMarshal(member))
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZRem(cc context.Context, key string, members ...T) *ldredis.IntCmd {
	return c.client.ZRem(cc, key, c.marshalSlice(members)...)
}
func (c *Redis[T]) ZRevRange(cc context.Context, key string, start, stop int64) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRevRange(cc, key, start, stop))
}
func (c *Redis[T]) ZRevRangeWithScores(cc context.Context, key string, start, stop int64) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRevRangeWithScores(cc, key, start, stop))
}
func (c *Redis[T]) ZRevRangeByScore(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRevRangeByScore(cc, key, opt))
}
func (c *Redis[T]) ZRevRangeByLex(cc context.Context, key string, opt *ldredis.ZRangeBy) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRevRangeByLex(cc, key, opt))
}
func (c *Redis[T]) ZRevRangeByScoreWithScores(cc context.Context, key string, opt *ldredis.ZRangeBy) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRevRangeByScoreWithScores(cc, key, opt))
}
func (c *Redis[T]) ZRevRank(cc context.Context, key string, member T) *ldredis.IntCmd {
	cmd := redis.NewIntCmd(cc, "zrevrank", key, c.mustMarshal(member))
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZRevRankWithScore(cc context.Context, key string, member T) *ldredis.RankWithScoreCmd {
	// return c.client.ZRevRankWithScore(cc, key, c.mustMarshal(member))
	cmd := redis.NewRankWithScoreCmd(cc, "zrevrank", key, c.mustMarshal(member), "withscore")
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZScore(cc context.Context, key string, member T) *ldredis.FloatCmd {
	cmd := redis.NewFloatCmd(cc, "zscore", key, c.mustMarshal(member))
	c.client.Process(cc, cmd)
	return cmd
}
func (c *Redis[T]) ZRandMember(cc context.Context, key string, count int) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZRandMember(cc, key, count))
}
func (c *Redis[T]) ZRandMemberWithScores(cc context.Context, key string, count int) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZRandMemberWithScores(cc, key, count))
}
func (c *Redis[T]) ZUnion(cc context.Context, store ldredis.ZStore) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZUnion(cc, store))
}
func (c *Redis[T]) ZUnionWithScores(cc context.Context, store ldredis.ZStore) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZUnionWithScores(cc, store))
}
func (c *Redis[T]) ZDiff(cc context.Context, keys ...string) *TypeSliceCmd[T] {
	return newTypeSliceCmd[T](cc, c.base, c.client.ZDiff(cc, keys...))
}
func (c *Redis[T]) ZDiffWithScores(cc context.Context, keys ...string) *ZMemberSliceCmd[T] {
	return newZMemberSliceCmd[T](cc, c.base, c.client.ZDiffWithScores(cc, keys...))
}

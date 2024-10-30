/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldredis"
	"github.com/smartystreets/goconvey/convey"
)

func testMemoryRedis() *ldredis.Redis {
	return ldredis.MustNewTestRedis()
}

func TestCodecRedis_Nil(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		type TestImpl struct{}

		var pTest *TestImpl
		var p interface{} = pTest

		c.So(pTest == nil, convey.ShouldBeTrue)
		c.So(p, convey.ShouldBeNil)
		c.So(p, convey.ShouldNotResemble, nil)

		key := "test-codec-redis-nil"
		expiration := 100 * time.Second

		s := New(rds, Json[any]()).SetNX(ctx, key, p, expiration)
		c.So(s.Err(), convey.ShouldBeNil)
		c.So(s.Val(), convey.ShouldBeTrue)

		g0 := rds.Get(ctx, key)
		c.So(g0.Err(), convey.ShouldBeNil)
		c.So(g0.Val(), convey.ShouldResemble, "")

		g1 := New(rds, Json[any]()).Get(ctx, key)
		c.So(g1.Err(), convey.ShouldBeNil)
		c.So(g1.Val(), convey.ShouldResemble, nil)
	})
}

func TestCodecRedis_String(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-codec-redis-string"
		expiration := time.Duration(0)

		// convey.Convey("marshal fail", func(c convey.C) {
		// 	cmd := New(rds, ProtoV1Codec[proto.Message]{}).Set(ctx, key, "234", expiration)
		// 	c.So(cmd.Err().Error(), convey.ShouldStartWith, "the object for marshal must be `proto.Message`.")
		// })
		c.Convey("set/setnx/setxx/get", func(c convey.C) {
			s0 := New(rds, Json[any]()).SetXX(ctx, key, "100", expiration)
			c.So(s0.Err(), convey.ShouldBeNil)
			c.So(s0.Val(), convey.ShouldBeFalse)

			s1 := New(rds, Json[any]()).SetNX(ctx, key, "100", expiration)
			c.So(s1.Err(), convey.ShouldBeNil)
			c.So(s1.Val(), convey.ShouldBeTrue)

			g := New(rds, Json[any]()).Get(ctx, key)
			c.So(g.Err(), convey.ShouldBeNil)
			c.So(g.Val(), convey.ShouldResemble, "100")

			g1 := rds.Get(ctx, key)
			c.So(g1.Err(), convey.ShouldBeNil)
			c.So(g1.Val(), convey.ShouldResemble, `"100"`)

			s2 := New(rds, Json[any]()).SetNX(ctx, key, "123", expiration)
			c.So(s2.Err(), convey.ShouldBeNil)
			c.So(s2.Val(), convey.ShouldBeFalse)

			s3 := New(rds, Json[any]()).SetXX(ctx, key, "123", expiration)
			c.So(s3.Err(), convey.ShouldBeNil)
			c.So(s3.Val(), convey.ShouldBeTrue)

			g3 := rds.Get(ctx, key)
			c.So(g3.Err(), convey.ShouldBeNil)
			c.So(g3.Val(), convey.ShouldResemble, `"123"`)

			s4 := New(rds, Json[any]()).Set(ctx, key, "234", expiration)
			c.So(s4.Err(), convey.ShouldBeNil)

			g4 := rds.Get(ctx, key)
			c.So(g4.Err(), convey.ShouldBeNil)
			c.So(g4.Val(), convey.ShouldResemble, `"234"`)
		})
		c.Convey("mset/msetnx/mget", func(c convey.C) {
			key0 := key + "-0"
			key1 := key + "-1"
			key2 := key + "-2"

			s0 := New(rds, Json[any]()).MSetNXPairs(ctx, Pair[any]{key0, "234"}, Pair[any]{key1, "abc"})
			c.So(s0.Err(), convey.ShouldBeNil)
			c.So(s0.Val(), convey.ShouldBeTrue)

			g := New(rds, Json[any]()).MGet(ctx, key0, key1, key2)
			c.So(g.Err(), convey.ShouldBeNil)
			c.So(g.Val(), convey.ShouldResemble, []interface{}{"234", "abc", nil})

			g0 := rds.MGet(ctx, key0, key1, key2)
			c.So(g0.Err(), convey.ShouldBeNil)
			c.So(g0.Val(), convey.ShouldResemble, []interface{}{`"234"`, `"abc"`, nil})

			s1 := New(rds, Json[any]()).MSetNXPairs(ctx, Pair[any]{key0, "100"}, Pair[any]{key2, "xyz"})
			c.So(s1.Err(), convey.ShouldBeNil)
			c.So(s1.Val(), convey.ShouldBeFalse)

			g1 := rds.MGet(ctx, key0, key1, key2)
			c.So(g1.Err(), convey.ShouldBeNil)
			c.So(g1.Val(), convey.ShouldResemble, []interface{}{`"234"`, `"abc"`, nil})

			s2 := New(rds, Json[any]()).MSetPairs(ctx, Pair[any]{key0, "100"}, Pair[any]{key2, "xyz"})
			c.So(s2.Err(), convey.ShouldBeNil)

			g2 := rds.MGet(ctx, key0, key1, key2)
			c.So(g2.Err(), convey.ShouldBeNil)
			c.So(g2.Val(), convey.ShouldResemble, []interface{}{`"100"`, `"abc"`, `"xyz"`})
		})
	})
}

func TestCodecRedis_Hash(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-codec-redis-hash"
		field := "field"

		c.Convey("hset/hsetnx/hget", func(c convey.C) {
			s0 := New(rds, Json[any]()).HSetNX(ctx, key, field, "100")
			c.So(s0.Err(), convey.ShouldBeNil)
			c.So(s0.Val(), convey.ShouldBeTrue)

			g0 := rds.HGet(ctx, key, field)
			c.So(g0.Err(), convey.ShouldBeNil)
			c.So(g0.Val(), convey.ShouldResemble, `"100"`)

			g := New(rds, Json[any]()).HGet(ctx, key, field)
			c.So(g.Err(), convey.ShouldBeNil)
			c.So(g.Val(), convey.ShouldResemble, "100")

			s1 := New(rds, Json[any]()).HSetNX(ctx, key, field, "100")
			c.So(s1.Err(), convey.ShouldBeNil)
			c.So(s1.Val(), convey.ShouldBeFalse)

			s2 := New(rds, Json[any]()).HSet(ctx, key, field, "abc")
			c.So(s2.Err(), convey.ShouldBeNil)

			g2 := rds.HGet(ctx, key, field)
			c.So(g2.Err(), convey.ShouldBeNil)
			c.So(g2.Val(), convey.ShouldResemble, `"abc"`)
		})
		c.Convey("hmset/hmget/hgetall/hvals", func(c convey.C) {
			field0 := field + "-0"
			field1 := field + "-1"
			field2 := field + "-2"

			s := New(rds, Json[any]()).HMSetMap(ctx, key, map[string]any{
				field0: "100",
				field1: 100,
			})
			c.So(s.Err(), convey.ShouldBeNil)

			g := rds.HMGet(ctx, key, field0, field1, field2)
			c.So(g.Err(), convey.ShouldBeNil)
			c.So(g.Val(), convey.ShouldResemble, []interface{}{`"100"`, "100", nil})

			g0 := New(rds, Json[any]()).HMGet(ctx, key, field0, field1, field2)
			c.So(g0.Err(), convey.ShouldBeNil)
			c.So(g0.Val(), convey.ShouldResemble, []interface{}{"100", float64(100), nil})

			g1 := New(rds, Json[any]()).HGetAll(ctx, key)
			c.So(g1.Err(), convey.ShouldBeNil)
			c.So(g1.Val(), convey.ShouldResemble, map[string]interface{}{field0: "100", field1: float64(100)})

			// d1 := New(rds, Json[any]()).HDel(ctx, key, field0)
			// c.So(d1.Err(), convey.ShouldBeNil)

			g2 := New(rds, Json[any]()).HVals(ctx, key)
			c.So(g2.Err(), convey.ShouldBeNil)
			c.So(g2.Val(), convey.ShouldHaveLength, 2)
			c.So(g2.Val(), convey.ShouldContain, float64(100))
			c.So(g2.Val(), convey.ShouldContain, "100")
		})
	})
}

func TestCodecRedis_ZSet(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-codec-redis-zset"

		members := []ZMember[any]{
			{Score: 1, Member: "100"},
			{Score: 10, Member: "abc"},
			{Score: 100, Member: "xyz"},
		}

		zadd := New(rds, Json[any]()).ZAdd(ctx, key, members...)
		c.So(zadd.Err(), convey.ShouldBeNil)

		c.Convey("zrange", func(c convey.C) {
			zrange0 := rds.ZRange(ctx, key, 0, -1)
			c.So(zrange0.Err(), convey.ShouldBeNil)
			c.So(zrange0.Val(), convey.ShouldResemble, []string{`"100"`, `"abc"`, `"xyz"`})

			zrange1 := New(rds, Json[any]()).ZRange(ctx, key, 0, -1)
			c.So(zrange1.Err(), convey.ShouldBeNil)
			c.So(zrange1.Val(), convey.ShouldResemble, []interface{}{"100", "abc", "xyz"})
		})

		c.Convey("zincrby/zscore", func(c convey.C) {
			zincrby := New(rds, Json[any]()).ZIncrBy(ctx, key, 10, "100")
			c.So(zincrby.Err(), convey.ShouldBeNil)
			c.So(zincrby.Val(), convey.ShouldResemble, float64(11))

			zscore := New(rds, Json[any]()).ZScore(ctx, key, "100")
			c.So(zscore.Err(), convey.ShouldBeNil)
			c.So(zscore.Val(), convey.ShouldResemble, float64(11))
		})

		c.Convey("zrank/zrevrank", func(c convey.C) {
			zrank0 := New(rds, Json[any]()).ZRank(ctx, key, "100")
			c.So(zrank0.Err(), convey.ShouldBeNil)
			c.So(zrank0.Val(), convey.ShouldResemble, int64(0))

			zrank1 := New(rds, Json[any]()).ZRank(ctx, key, "aaa")
			c.So(zrank1.Err(), convey.ShouldEqual, ldredis.Nil)

			zrevrank := New(rds, Json[any]()).ZRevRank(ctx, key, "100")
			c.So(zrevrank.Err(), convey.ShouldBeNil)
			c.So(zrevrank.Val(), convey.ShouldResemble, int64(2))
		})

		c.Convey("zrem", func(c convey.C) {
			zrem := New(rds, Json[any]()).ZRem(ctx, key, "abc")
			c.So(zrem.Err(), convey.ShouldBeNil)
			c.So(zrem.Val(), convey.ShouldResemble, int64(1))

			zrange := rds.ZRange(ctx, key, 0, -1)
			c.So(zrange.Err(), convey.ShouldBeNil)
			c.So(zrange.Val(), convey.ShouldResemble, []string{`"100"`, `"xyz"`})
		})
	})
}

func TestCodecRedis_Set(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key0 := "test-codec-redis-set-0"
		key1 := "test-codec-redis-set-1"

		sadd0 := New(rds, Json[any]()).SAdd(ctx, key0, "abc", 100, "100")
		c.So(sadd0.Err(), convey.ShouldBeNil)

		sadd1 := New(rds, Json[any]()).SAdd(ctx, key1, "abc", 234, "xyz")
		c.So(sadd1.Err(), convey.ShouldBeNil)

		c.Convey("smembers", func(c convey.C) {
			g0 := rds.SMembers(ctx, key0)
			c.So(g0.Err(), convey.ShouldBeNil)
			c.So(g0.Val(), convey.ShouldHaveLength, 3)
			c.So(g0.Val(), convey.ShouldContain, `"100"`)
			c.So(g0.Val(), convey.ShouldContain, `100`)
			c.So(g0.Val(), convey.ShouldContain, `"abc"`)

			g1 := New(rds, Json[any]()).SMembers(ctx, key0)
			c.So(g1.Err(), convey.ShouldBeNil)
			c.So(g1.Val(), convey.ShouldHaveLength, 3)
			c.So(g1.Val(), convey.ShouldContain, "100")
			c.So(g1.Val(), convey.ShouldContain, float64(100))
			c.So(g1.Val(), convey.ShouldContain, "abc")
		})
		c.Convey("smembers map", func(c convey.C) {
			g := New(rds, Json[any]()).SMembersMap(ctx, key0)
			c.So(g.Err(), convey.ShouldBeNil)
			c.So(g.Val(), convey.ShouldResemble, map[interface{}]struct{}{
				"abc":        {},
				"100":        {},
				float64(100): {},
			})
		})
		c.Convey("sismember", func(c convey.C) {
			g0 := New(rds, Json[any]()).SIsMember(ctx, key0, 100)
			c.So(g0.Err(), convey.ShouldBeNil)
			c.So(g0.Val(), convey.ShouldBeTrue)

			g1 := New(rds, Json[any]()).SIsMember(ctx, key0, "xyz")
			c.So(g1.Err(), convey.ShouldBeNil)
			c.So(g1.Val(), convey.ShouldBeFalse)
		})
		c.Convey("sdiff/sinter", func(c convey.C) {
			sinter := New(rds, Json[any]()).SInter(ctx, key0, key1)
			c.So(sinter.Err(), convey.ShouldBeNil)
			c.So(sinter.Val(), convey.ShouldResemble, []interface{}{"abc"})

			sdiff := New(rds, Json[any]()).SDiff(ctx, key0, key1)
			c.So(sdiff.Err(), convey.ShouldBeNil)
			c.So(sdiff.Val(), convey.ShouldHaveLength, 2)
			c.So(sdiff.Val(), convey.ShouldContain, float64(100))
			c.So(sdiff.Val(), convey.ShouldContain, "100")
		})
		c.Convey("srem", func(c convey.C) {
			srem := New(rds, Json[any]()).SRem(ctx, key0, "100")
			c.So(srem.Err(), convey.ShouldBeNil)
			c.So(srem.Val(), convey.ShouldEqual, 1)

			smembers := New(rds, Json[any]()).SMembers(ctx, key0)
			c.So(smembers.Err(), convey.ShouldBeNil)
			c.So(smembers.Val(), convey.ShouldHaveLength, 2)
			c.So(smembers.Val(), convey.ShouldContain, float64(100))
			c.So(smembers.Val(), convey.ShouldContain, "abc")
		})
	})
}

func TestCodecRedis_List(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := testMemoryRedis()
		defer rds.Close()

		ctx := ldctx.Discard()

		key := "test-codec-redis-list"

		lpush := New(rds, Json[any]()).LPush(ctx, key, "abc", 100, "100")
		c.So(lpush.Err(), convey.ShouldBeNil)
		c.So(lpush.Val(), convey.ShouldEqual, 3)

		rpush := New(rds, Json[any]()).RPush(ctx, key, "xyz", 234, "100.234")
		c.So(rpush.Err(), convey.ShouldBeNil)
		c.So(rpush.Val(), convey.ShouldEqual, 6)

		c.Convey("llen", func(c convey.C) {
			cmd := New(rds, Json[any]()).LLen(ctx, key)
			c.So(cmd.Err(), convey.ShouldBeNil)
			c.So(cmd.Val(), convey.ShouldEqual, 6)
		})
		c.Convey("lrange", func(c convey.C) {
			cmd := New(rds, Json[any]()).LRange(ctx, key, 0, -1)
			c.So(cmd.Err(), convey.ShouldBeNil)
			c.So(cmd.Val(), convey.ShouldResemble, []interface{}{
				"100", float64(100), "abc", "xyz", float64(234), "100.234",
			})
		})
		c.Convey("lindex", func(c convey.C) {
			cmd0 := New(rds, Json[any]()).LIndex(ctx, key, 0)
			c.So(cmd0.Err(), convey.ShouldBeNil)
			c.So(cmd0.Val(), convey.ShouldResemble, "100")

			cmd1 := New(rds, Json[any]()).LIndex(ctx, key, 5)
			c.So(cmd1.Err(), convey.ShouldBeNil)
			c.So(cmd1.Val(), convey.ShouldResemble, "100.234")
		})
		c.Convey("lpop", func(c convey.C) {
			cmd := New(rds, Json[any]()).LPop(ctx, key)
			c.So(cmd.Err(), convey.ShouldBeNil)
			c.So(cmd.Val(), convey.ShouldResemble, "100")
		})
		c.Convey("rpop", func(c convey.C) {
			cmd := New(rds, Json[any]()).RPop(ctx, key)
			c.So(cmd.Err(), convey.ShouldBeNil)
			c.So(cmd.Val(), convey.ShouldResemble, "100.234")
		})
		c.Convey("LRem", func(c convey.C) {
			lpush := New(rds, Json[any]()).LPush(ctx, key, "100", "aaa", "100", "bbb")
			c.So(lpush.Err(), convey.ShouldBeNil)
			c.So(lpush.Val(), convey.ShouldEqual, 10)

			lrem := New(rds, Json[any]()).LRem(ctx, key, 2, "100")
			c.So(lrem.Err(), convey.ShouldBeNil)
			c.So(lrem.Val(), convey.ShouldEqual, 2)

			lrange := New(rds, Json[any]()).LRange(ctx, key, 0, -1)
			c.So(lrange.Err(), convey.ShouldBeNil)
			c.So(lrange.Val(), convey.ShouldResemble, []interface{}{
				"bbb", "aaa", "100", float64(100), "abc", "xyz", float64(234), "100.234",
			})
		})
	})
}

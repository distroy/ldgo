/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"time"

	"github.com/distroy/ldgo/v3/ldredis/internal"
	redis "github.com/redis/go-redis/v9"
)

type (
	Reporter = internal.Reporter
)

type (
	Cmd   = redis.Cmd
	Cmder = redis.Cmder

	PubSub       = redis.PubSub
	Message      = redis.Message
	Subscription = redis.Subscription

	Sort         = redis.Sort
	ScanCmd      = redis.ScanCmd
	ScanIterator = redis.ScanIterator

	SetArgs  = redis.SetArgs
	LPosArgs = redis.LPosArgs

	BoolCmd   = redis.BoolCmd
	StatusCmd = redis.StatusCmd

	DurationCmd = redis.DurationCmd
	StringCmd   = redis.StringCmd
	IntCmd      = redis.IntCmd
	FloatCmd    = redis.FloatCmd

	SliceCmd           = redis.SliceCmd
	IntSliceCmd        = redis.IntSliceCmd
	BoolSliceCmd       = redis.BoolSliceCmd
	FloatSliceCmd      = redis.FloatSliceCmd
	StringSliceCmd     = redis.StringSliceCmd
	MapStringStringCmd = redis.MapStringStringCmd
	StringSetCmd       = redis.StringStructMapCmd

	KeyValue         = redis.KeyValue
	KeyValuesCmd     = redis.KeyValuesCmd
	KeyValueSliceCmd = redis.KeyValueSliceCmd

	RankScore        = redis.RankScore
	RankWithScoreCmd = redis.RankWithScoreCmd

	ZMember          = redis.Z
	ZAddArgs         = redis.ZAddArgs
	ZStore           = redis.ZStore
	ZRangeArgs       = redis.ZRangeArgs
	ZRangeBy         = redis.ZRangeBy
	ZSliceCmd        = redis.ZSliceCmd
	ZWithKey         = redis.ZWithKey
	ZWithKeyCmd      = redis.ZWithKeyCmd
	ZSliceWithKeyCmd = redis.ZSliceWithKeyCmd

	XAddArgs         = redis.XAddArgs
	XReadArgs        = redis.XReadArgs
	XReadGroupArgs   = redis.XReadGroupArgs
	XClaimArgs       = redis.XClaimArgs
	XPending         = redis.XPending
	XPendingCmd      = redis.XPendingCmd
	XPendingExt      = redis.XPendingExt
	XPendingExtArgs  = redis.XPendingExtArgs
	XPendingExtCmd   = redis.XPendingExtCmd
	XMessage         = redis.XMessage
	XMessageSliceCmd = redis.XMessageSliceCmd
	XStream          = redis.XStream
	XStreamSliceCmd  = redis.XStreamSliceCmd

	GeoLocation    = redis.GeoLocation
	GeoLocationCmd = redis.GeoLocationCmd
	GeoPos         = redis.GeoPos
	GeoPosCmd      = redis.GeoPosCmd
	GeoRadiusQuery = redis.GeoRadiusQuery
)

type StringCmdable interface {
	redis.StringCmdable

	// MSet is like Set but accepts multiple values:
	//   - MSet("key1", "value1", "key2", "value2")
	//   - MSet([]string{"key1", "value1", "key2", "value2"})
	//   - MSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
	//   - MSet(struct), For struct types, see HSet description.
	MSet(ctx context.Context, values ...interface{}) *StatusCmd
	// MSetNX is like SetNX but accepts multiple values:
	//   - MSetNX("key1", "value1", "key2", "value2")
	//   - MSetNX([]string{"key1", "value1", "key2", "value2"})
	//   - MSetNX(map[string]interface{}{"key1": "value1", "key2": "value2"})
	//   - MSetNX(struct), For struct types, see HSet description.
	MSetNX(ctx context.Context, values ...interface{}) *BoolCmd
	// Set Redis `SET key value [expiration]` command.
	// Use expiration for `SETEx`-like behavior.
	//
	// Zero expiration means the key has no expiration time.
	// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
	// otherwise you will receive an error: (error) ERR syntax error.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
}

type HashCmdable interface {
	redis.HashCmdable

	// HSet accepts values in following formats:
	//
	//   - HSet("myhash", "key1", "value1", "key2", "value2")
	//
	//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
	//
	//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
	//
	//     Playing struct With "redis" tag.
	//     type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
	//
	//   - HSet("myhash", MyHash{"value1", "value2"}) Warn: redis-server >= 4.0
	//
	//     For struct, can be a structure pointer type, we only parse the field whose tag is redis.
	//     if you don't want the field to be read, you can use the `redis:"-"` flag to ignore it,
	//     or you don't need to set the redis tag.
	//     For the type of structure field, we only support simple data types:
	//     string, int/uint(8,16,32,64), float(32,64), time.Time(to RFC3339Nano), time.Duration(to Nanoseconds ),
	//     if you are other more complex or custom data types, please implement the encoding.BinaryMarshaler interface.
	//
	// Note that in older versions of Redis server(redis-server < 4.0), HSet only supports a single key-value pair.
	// redis-docs: https://redis.io/commands/hset (Starting with Redis version 4.0.0: Accepts multiple field and value arguments.)
	// If you are using a Struct type and the number of fields is greater than one,
	// you will receive an error similar to "ERR wrong number of arguments", you can use HMSet as a substitute.
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd

	// HMSet is a deprecated version of HSet left for compatibility with Redis 3.
	HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd
}

type SortedSetCmdable interface {
	redis.SortedSetCmdable

	// BZMPop is the blocking variant of ZMPOP.
	// When any of the sorted sets contains elements, this command behaves exactly like ZMPOP.
	// When all sorted sets are empty, Redis will block the connection until another client adds members to one of the keys or until the timeout elapses.
	// A timeout of zero can be used to block indefinitely.
	// example: client.BZMPop(ctx, 0,"max", 1, "set")
	BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZSliceWithKeyCmd
}

type Cmdable interface {
	StringCmdable
	HashCmdable
	SortedSetCmdable

	redis.Cmdable

	Do(ctx context.Context, args ...interface{}) *Cmd
	Process(ctx context.Context, cmd Cmder) error
	Close() error

	Subscribe(ctx context.Context, channels ...string) *PubSub
	PSubscribe(ctx context.Context, channels ...string) *PubSub
}

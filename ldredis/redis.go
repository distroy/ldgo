/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"time"

	"github.com/distroy/ldgo/v3/ldredis/internal"
	redis "github.com/redis/go-redis/v9"
)

const (
	Nil = redis.Nil
)

func isErrNil(err error) bool {
	return err == nil || err == Nil
}

var _ Cmdable = (*Redis)(nil)

const (
	redisRetryIntervalMin = 10 * time.Millisecond
	redisRetryIntervalDef = 1 * time.Second
)

func New(cli redis.Cmdable) *Redis {
	switch v := cli.(type) {
	case *Redis:
		return v

	case wrapper:
		return newRedis(v)

	case *redis.Client:
		return newRedis(newWrapper(v))

	case *redis.ClusterClient:
		return newRedis(newWrapper(v))

	case cmdable:
		return newRedis(newWrapper(v))
	}

	panic("redis client type must be `*ldredis.Redis` or `*redis.Client` or `*redis.ClusterClient`")
}

func NewByConfig(cfg *Config) *Redis {
	if cfg.Cluster {
		return newRedis(newRedisCluster(cfg))
	}

	return newRedis(newRedisClient(cfg))
}

func newRedis(cli wrapper) *Redis {
	c := &Redis{
		origin: cli,
		opts: internal.Options{
			Reporter: internal.DiscardReporter{},
			Caller:   true,
		},
	}

	c.wrapper = c.origin.Clone()
	c.wrapper.AddHook(newHook(c))
	return c
}

// Redis struct
type Redis struct {
	wrapper

	origin wrapper
	opts   internal.Options
}

func (c *Redis) Client() Cmdable { return c.origin.cmdable }
func (c *Redis) Clone() *Redis   { return c.clone() }

func (c *Redis) clone() *Redis {
	cp := *c
	c = &cp
	c.wrapper = c.origin.Clone()
	c.wrapper.AddHook(newHook(c))
	return c
}

// WithRetry should be called like these:
//
//	WithRetry(3)
//	WithRetry(3, {interval})
func (c *Redis) WithRetry(retry int, interval ...time.Duration) *Redis {
	c = c.clone()
	c.opts.Retry = retry

	c.opts.RetryInterval = redisRetryIntervalDef
	if len(interval) > 0 {
		d := interval[0]
		if d < redisRetryIntervalMin {
			d = redisRetryIntervalMin
		}
		c.opts.RetryInterval = d
	}

	return c
}

func (c *Redis) WithReport(reporter Reporter) *Redis {
	if reporter == nil {
		reporter = internal.DiscardReporter{}
	}

	c = c.clone()
	c.opts.Reporter = reporter
	return c
}

func (c *Redis) WithCaller(enable bool) *Redis {
	c = c.clone()
	c.opts.Caller = enable
	return c
}

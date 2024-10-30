/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"github.com/distroy/ldgo/v3/ldredis"
)

func New[T comparable](rds *ldredis.Redis, codec Codec[T]) *Redis[T] {
	return &Redis[T]{
		base: base[T]{
			client: rds,
			codec:  codec,
		},
	}
}

var (
	_ Cmdable[int]    = (*Redis[int])(nil)
	_ Cmdable[string] = (*Redis[string])(nil)
)

type Redis[T comparable] struct {
	base[T]
}

func (c *Redis[T]) Client() *ldredis.Redis { return c.client }
func (c *Redis[T]) Codec() Codec[T]        { return c.codec }

func (c *Redis[T]) clone(cli ...*ldredis.Redis) *Redis[T] {
	cp := *c
	c = &cp
	if len(cli) > 0 && cli[0] != nil {
		cp.client = cli[0]
	}
	return c
}

func (c *Redis[T]) WithRetry(retry int) *Redis[T] {
	return c.clone(c.client.WithRetry(retry))
}
func (c *Redis[T]) WithReport(r Reporter) *Redis[T] {
	return c.clone(c.client.WithReport(r))
}
func (c *Redis[T]) WithCaller(enable bool) *Redis[T] {
	return c.clone(c.client.WithCaller(enable))
}

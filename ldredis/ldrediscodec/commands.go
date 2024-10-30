/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"context"
	"encoding"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldredis"
	"github.com/distroy/ldgo/v3/ldredis/internal"
	redis "github.com/redis/go-redis/v9"
)

var (
	_ encoding.BinaryMarshaler = (*errorMarshaler)(nil)

	_ internal.CmderWithParse = (*TypeCmd[any])(nil)
	_ internal.CmderWithParse = (*MapStringTypeCmd[any])(nil)
	_ internal.CmderWithParse = (*TypeSliceCmd[any])(nil)
	_ internal.CmderWithParse = (*SliceCmd[any])(nil)
	_ internal.CmderWithParse = (*ZMemberSliceCmd[any])(nil)
)

type errorMarshaler struct {
	Err error
}

func (c errorMarshaler) MarshalBinary() ([]byte, error) {
	return nil, c.Err
}

func newTypeCmd[T any](cc context.Context, base base[T], cmd *ldredis.StringCmd) *TypeCmd[T] {
	c := &TypeCmd[T]{
		base:      base,
		StringCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewTypeCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *TypeCmd[T] {
	c := &TypeCmd[T]{
		base:      base[T]{codec: codec},
		StringCmd: redis.NewStringCmd(cc, args...),
	}
	return c
}

type TypeCmd[T any] struct {
	base[T]
	*ldredis.StringCmd

	val T
}

func (c *TypeCmd[T]) Val() T             { return c.val }
func (c *TypeCmd[T]) Result() (T, error) { return c.Val(), c.Err() }
func (c *TypeCmd[T]) Parse(cc context.Context) error {
	cmd := c.StringCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	v, err := c.unmarshal(cc, ldconv.StrToBytesUnsafe(cmd.Val()))
	if err != nil {
		c.SetErr(err)
		return err
	}

	c.val = v
	return nil
}

func newMapStringTypeCmd[T any](cc context.Context, base base[T], cmd *ldredis.MapStringStringCmd) *MapStringTypeCmd[T] {
	c := &MapStringTypeCmd[T]{
		base:               base,
		MapStringStringCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewMapStringTypeCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *MapStringTypeCmd[T] {
	c := &MapStringTypeCmd[T]{
		base:               base[T]{codec: codec},
		MapStringStringCmd: redis.NewMapStringStringCmd(cc, args...),
	}
	return c
}

type MapStringTypeCmd[T any] struct {
	base[T]
	*ldredis.MapStringStringCmd

	val map[string]T
}

func (c *MapStringTypeCmd[T]) Val() map[string]T             { return c.val }
func (c *MapStringTypeCmd[T]) Result() (map[string]T, error) { return c.Val(), c.Err() }
func (c *MapStringTypeCmd[T]) Parse(cc context.Context) error {
	cmd := c.MapStringStringCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	m := make(map[string]T, len(cmd.Val()))
	for k, val := range cmd.Val() {
		v, err := c.unmarshal(cc, ldconv.StrToBytesUnsafe(val))
		if err != nil {
			c.SetErr(err)
			return err
		}
		m[k] = v
	}
	c.val = m
	return nil
}

func newTypeSliceCmd[T any](cc context.Context, base base[T], cmd *ldredis.StringSliceCmd) *TypeSliceCmd[T] {
	c := &TypeSliceCmd[T]{
		base:           base,
		StringSliceCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewTypeSliceCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *TypeSliceCmd[T] {
	c := &TypeSliceCmd[T]{
		base:           base[T]{codec: codec},
		StringSliceCmd: redis.NewStringSliceCmd(cc, args...),
	}
	return c
}

type TypeSliceCmd[T any] struct {
	base[T]
	*ldredis.StringSliceCmd

	val []T
}

func (c *TypeSliceCmd[T]) Val() []T             { return c.val }
func (c *TypeSliceCmd[T]) Result() ([]T, error) { return c.Val(), c.Err() }
func (c *TypeSliceCmd[T]) Parse(cc context.Context) error {
	cmd := c.StringSliceCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	s := make([]T, 0, len(cmd.Val()))
	for _, val := range cmd.Val() {
		v, err := c.unmarshal(cc, ldconv.StrToBytesUnsafe(val))
		if err != nil {
			c.SetErr(err)
			return err
		}
		s = append(s, v)
	}
	c.val = s
	return nil
}

func newTypeSetCmd[T comparable](cc context.Context, base base[T], cmd *ldredis.StringSetCmd) *TypeSetCmd[T] {
	c := &TypeSetCmd[T]{
		base:         base,
		StringSetCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewTypeSetCmd[T comparable](cc context.Context, codec Codec[T], args ...interface{}) *TypeSetCmd[T] {
	c := &TypeSetCmd[T]{
		base:         base[T]{codec: codec},
		StringSetCmd: redis.NewStringStructMapCmd(cc, args...),
	}
	return c
}

type TypeSetCmd[T comparable] struct {
	base[T]
	*ldredis.StringSetCmd

	val map[T]struct{}
}

func (c *TypeSetCmd[T]) Val() map[T]struct{}             { return c.val }
func (c *TypeSetCmd[T]) Result() (map[T]struct{}, error) { return c.Val(), c.Err() }
func (c *TypeSetCmd[T]) Parse(cc context.Context) error {
	cmd := c.StringSetCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	s := make(map[T]struct{}, len(cmd.Val()))
	for val := range cmd.Val() {
		v, err := c.unmarshal(cc, ldconv.StrToBytesUnsafe(val))
		if err != nil {
			c.SetErr(err)
			return err
		}
		s[v] = struct{}{}
	}
	c.val = s
	return nil
}

func newSliceCmd[T any](cc context.Context, base base[T], cmd *ldredis.SliceCmd) *SliceCmd[T] {
	c := &SliceCmd[T]{
		base:     base,
		SliceCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewSliceCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *SliceCmd[T] {
	c := &SliceCmd[T]{
		base:     base[T]{codec: codec},
		SliceCmd: redis.NewSliceCmd(cc, args...),
	}
	return c
}

type SliceCmd[T any] struct {
	base[T]
	*ldredis.SliceCmd

	val []T
}

func (c *SliceCmd[T]) Val() []T             { return c.val }
func (c *SliceCmd[T]) Result() ([]T, error) { return c.Val(), c.Err() }
func (c *SliceCmd[T]) Parse(cc context.Context) error {
	cmd := c.SliceCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	s := make([]T, 0, len(cmd.Val()))
	for _, val := range cmd.Val() {
		v, err := c.unmarshalInterface(cc, val)
		if err != nil {
			c.SetErr(err)
			return err
		}
		s = append(s, v)
	}
	c.val = s
	return nil
}

func newZMemberSliceCmd[T any](cc context.Context, base base[T], cmd *ldredis.ZSliceCmd) *ZMemberSliceCmd[T] {
	c := &ZMemberSliceCmd[T]{
		base:      base,
		ZSliceCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewZMemberSliceCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *ZMemberSliceCmd[T] {
	c := &ZMemberSliceCmd[T]{
		base:      base[T]{codec: codec},
		ZSliceCmd: redis.NewZSliceCmd(cc, args...),
	}
	return c
}

type ZMemberSliceCmd[T any] struct {
	base[T]
	*ldredis.ZSliceCmd

	val []ZMember[T]
}

func (c *ZMemberSliceCmd[T]) Val() []ZMember[T]             { return c.val }
func (c *ZMemberSliceCmd[T]) Result() ([]ZMember[T], error) { return c.Val(), c.Err() }
func (c *ZMemberSliceCmd[T]) Parse(cc context.Context) error {
	cmd := c.ZSliceCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	members := make([]ZMember[T], 0, len(cmd.Val()))
	for _, v := range cmd.Val() {
		val, err := c.unmarshalInterface(cc, v.Member)
		if err != nil {
			c.SetErr(err)
			return err
		}
		members = append(members, ZMember[T]{
			Score:  v.Score,
			Member: val,
		})
	}

	c.val = members
	return nil
}

func newZMemberWithKeyCmd[T any](cc context.Context, base base[T], cmd *ldredis.ZWithKeyCmd) *ZMemberWithKeyCmd[T] {
	c := &ZMemberWithKeyCmd[T]{
		base:        base,
		ZWithKeyCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewZMemberWithKeyCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *ZMemberWithKeyCmd[T] {
	c := &ZMemberWithKeyCmd[T]{
		base:        base[T]{codec: codec},
		ZWithKeyCmd: redis.NewZWithKeyCmd(cc, args...),
	}
	return c
}

type ZMemberWithKeyCmd[T any] struct {
	base[T]
	*ldredis.ZWithKeyCmd

	val ZMemberWithKey[T]
}

func (c *ZMemberWithKeyCmd[T]) Val() ZMemberWithKey[T]             { return c.val }
func (c *ZMemberWithKeyCmd[T]) Result() (ZMemberWithKey[T], error) { return c.Val(), c.Err() }
func (c *ZMemberWithKeyCmd[T]) Parse(cc context.Context) error {
	cmd := c.ZWithKeyCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	v := cmd.Val()
	val, err := c.unmarshalInterface(cc, v.Member)
	if err != nil {
		c.SetErr(err)
		return err
	}

	c.val = ZMemberWithKey[T]{
		ZMember: ZMember[T]{
			Score:  v.Score,
			Member: val,
		},
		Key: v.Key,
	}

	return nil
}

func newZMemberSliceWithKeyCmd[T any](cc context.Context, base base[T], cmd *ldredis.ZSliceWithKeyCmd) *ZMemberSliceWithKeyCmd[T] {
	c := &ZMemberSliceWithKeyCmd[T]{
		base:             base,
		ZSliceWithKeyCmd: cmd,
	}
	c.Parse(cc)
	return c
}

func NewZMemberSliceWithKeyCmd[T any](cc context.Context, codec Codec[T], args ...interface{}) *ZMemberSliceWithKeyCmd[T] {
	c := &ZMemberSliceWithKeyCmd[T]{
		base:             base[T]{codec: codec},
		ZSliceWithKeyCmd: redis.NewZSliceWithKeyCmd(cc, args...),
	}
	return c
}

type ZMemberSliceWithKeyCmd[T any] struct {
	base[T]
	*ldredis.ZSliceWithKeyCmd

	key string
	val []ZMember[T]
}

func (c *ZMemberSliceWithKeyCmd[T]) Val() (string, []ZMember[T]) { return c.key, c.val }
func (c *ZMemberSliceWithKeyCmd[T]) Result() (string, []ZMember[T], error) {
	return c.key, c.val, c.Err()
}
func (c *ZMemberSliceWithKeyCmd[T]) Parse(cc context.Context) error {
	cmd := c.ZSliceWithKeyCmd
	if err := cmd.Err(); err != nil {
		return err
	}

	key, v := cmd.Val()
	members := make([]ZMember[T], 0, len(v))
	for _, v := range v {
		val, err := c.unmarshalInterface(cc, v.Member)
		if err != nil {
			c.SetErr(err)
			return err
		}
		members = append(members, ZMember[T]{
			Score:  v.Score,
			Member: val,
		})
	}
	c.key = key
	c.val = members
	return nil
}

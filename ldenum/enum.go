/*
 * Copyright (C) distroy
 */

package ldenum

import "github.com/distroy/ldgo/v3/ldptr"

type Enum[T interface{ EnumToString(n int) string }] int

func (n Enum[T]) Ptr() *Enum[T] { return &n }
func (n Enum[T]) Int() int      { return int(n) }
func (n Enum[T]) Str() string {
	var x T
	return (x).EnumToString(n.Int())
}

func (n *Enum[T]) Get() Enum[T]   { return ldptr.Get(n) }
func (n *Enum[T]) GetInt() int    { return n.Get().Int() }
func (n *Enum[T]) GetStr() string { return n.Get().Str() }

func (n *Enum[T]) New() *Enum[T] { return ldptr.NewByPtr(n) }
func (n *Enum[T]) NewInt() *int  { return (*int)(n.New()) }
func (n *Enum[T]) NewStr() *string {
	if n == nil {
		return nil
	}
	return ldptr.New(n.Get().Str())
}

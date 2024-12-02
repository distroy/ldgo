/*
 * Copyright (C) distroy
 */

package ldenum

import (
	"fmt"
	"sync"

	"github.com/distroy/ldgo/v2/ldptr"
)

var enumStringMap = &sync.Map{}

func NewEnumString[T ~int](name string, m map[T]string) *EnumString[T] {
	var k *EnumString[T]
	p := &EnumString[T]{
		Name: name,
		Map:  m,
	}

	enumStringMap.Store(k, p)
	return p
}

type EnumString[T ~int] struct {
	Name string
	Map  map[T]string
}

func (_ *EnumString[T]) EnumToString(n int) string {
	var k *EnumString[T]
	i, _ := enumStringMap.Load(k)
	x := i.(*EnumString[T])
	s := x.Map[T(n)]
	if s != "" {
		return s
	}
	return fmt.Sprintf("%s[%d]", x.Name, n)
}

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

/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"unsafe"
)

type iterator[T any] struct {
	ptr *T
	idx int
}

type Iterator[T any] iterator[T]

func (i Iterator[T]) Set(d T)                { unsafe.Slice(i.ptr, i.idx+1)[i.idx] = d }
func (i Iterator[T]) Get() T                 { return unsafe.Slice(i.ptr, i.idx+1)[i.idx] }
func (i Iterator[T]) Prev() Iterator[T]      { return i.Move(-1) }
func (i Iterator[T]) Next() Iterator[T]      { return i.Move(+1) }
func (i Iterator[T]) Move(n int) Iterator[T] { return Iterator[T]{ptr: i.ptr, idx: i.idx + n} }

type ReverseIterator[T any] iterator[T]

func (i ReverseIterator[T]) Set(d T)                  { unsafe.Slice(i.ptr, i.idx+1)[i.idx] = d }
func (i ReverseIterator[T]) Get() T                   { return unsafe.Slice(i.ptr, i.idx+1)[i.idx] }
func (i ReverseIterator[T]) Prev() ReverseIterator[T] { return i.Move(+1) }
func (i ReverseIterator[T]) Next() ReverseIterator[T] { return i.Move(-1) }
func (i ReverseIterator[T]) Move(n int) ReverseIterator[T] {
	return ReverseIterator[T]{ptr: i.ptr, idx: i.idx + n}
}

func makeIterator[T any](s []T, idx int) iterator[T] {
	return iterator[T]{ptr: unsafe.SliceData(s), idx: idx}
}

func Begin[T any](s []T) Iterator[T] { return Iterator[T](makeIterator(s, 0)) }
func End[T any](s []T) Iterator[T]   { return Iterator[T](makeIterator(s, len(s))) }

func RBegin[T any](s []T) ReverseIterator[T] { return ReverseIterator[T](makeIterator(s, len(s)-1)) }
func REnd[T any](s []T) ReverseIterator[T]   { return ReverseIterator[T](makeIterator(s, -1)) }

type Range[T any] struct {
	Begin Iterator[T]
	End   Iterator[T]
}

func (r *Range[T]) Set(d T)       { r.Begin.Set(d) }
func (r *Range[T]) Get() T        { return r.Begin.Get() }
func (r *Range[T]) HasNext() bool { return r.Begin.ptr != nil && r.Begin != r.End }
func (r *Range[T]) Next()         { r.Begin = r.Begin.Next() }

type ReverseRange[T any] struct {
	Begin ReverseIterator[T]
	End   ReverseIterator[T]
}

func (r *ReverseRange[T]) Set(d T)       { r.Begin.Set(d) }
func (r *ReverseRange[T]) Get() T        { return r.Begin.Get() }
func (r *ReverseRange[T]) HasNext() bool { return r.Begin.ptr != nil && r.Begin != r.End }
func (r *ReverseRange[T]) Next()         { r.Begin = r.Begin.Next() }

/*
 * Copyright (C) distroy
 */

package ldatomic

var (
	_ Loader[any] = funcLoader[any](nil)
	_ Storer[any] = funcStorer[any](nil)

	_ StoreLoadCompareAndSwapper[any] = (*Any[any])(nil)
)

func NewLoader[T any](f func() T) Loader[T]  { return funcLoader[T](f) }
func NewStorer[T any](f func(d T)) Storer[T] { return funcStorer[T](f) }

type Loader[T any] interface {
	Load() T
}

type Storer[T any] interface {
	Store(d T)
}

type StoreLoader[T any] interface {
	Loader[T]
	Storer[T]
}

type StoreLoadSwapper[T any] interface {
	StoreLoader[T]

	Swap(new T) (old T)
}

type StoreLoadCompareAndSwapper[T any] interface {
	StoreLoader[T]

	Swap(new T) (old T)
	CompareAndSwap(old, new T) (swapped bool)
}

type funcLoader[T any] func() T

func (f funcLoader[T]) Load() T { return f() }

type funcStorer[T any] func(d T)

func (f funcStorer[T]) Store(d T) { f(d) }

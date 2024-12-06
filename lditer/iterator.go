/*
 * Copyright (C) distroy
 */

package lditer

import "iter"

func Chan[T any](ch <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range ch {
			if !yield(v) {
				break
			}
		}
	}
}

func Chan2[T any](ch <-chan T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range ch {
			if !yield(i, v) {
				break
			}
			i++
		}
	}
}

func Slice[T any](s []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				break
			}
		}
	}
}

func Slice2[T any](s []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range s {
			if !yield(i, v) {
				break
			}
		}
	}
}

func Map[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func MapKeys[K comparable, V any](m map[K]V) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				break
			}
		}
	}
}

func MapValues[K comparable, V any](m map[K]V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				break
			}
		}
	}
}

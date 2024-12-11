/*
 * Copyright (C) distroy
 */

package lditer

import (
	"iter"
)

func ToSeq2[Seq ~func(yield func(V) bool), V any](fn Seq) iter.Seq2[int, V] {
	return func(yield func(i int, v V) bool) {
		x := 0
		fn(func(vv V) bool {
			i := x
			x++
			return yield(i, vv)
		})
	}
}

func ToSeqByKey[Seq2 ~func(yield func(k K, v V) bool), K, V any](fn Seq2) iter.Seq[K] {
	return func(yield func(K) bool) {
		fn(func(k K, v V) bool { return yield(k) })
	}
}

func ToSeqByValue[Seq2 ~func(yield func(k K, v V) bool), K, V any](fn Seq2) iter.Seq[V] {
	return func(yield func(V) bool) {
		fn(func(k K, v V) bool { return yield(v) })
	}
}

func Chan[C interface{ ~<-chan V | ~chan V }, V any](ch C) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range ch {
			if !yield(v) {
				break
			}
		}
	}
}

func Slice[S ~[]V, V any](s S) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range s {
			if !yield(i, v) {
				break
			}
		}
	}
}

func SliceBackward[S ~[]V, V any](s S) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

func Map[M ~map[K]V, K comparable, V any](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}
	}
}

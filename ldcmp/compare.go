/*
 * Copyright (C) distroy
 */

package ldcmp

import (
	"github.com/distroy/ldgo/v3/internal/cmp"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

type Complex interface {
	~complex64 | ~complex128
}

type Comparable interface {
	Number | ~string
}

type Comparer[T any] interface {
	Compare(b T) int
}

func Compare(a, b interface{}) int { return cmp.Compare(a, b) }

// Deprecated: use `Compare` instead.
func CompareInterface(a, b interface{}) int { return cmp.Compare(a, b) }

func CompareBool[T ~bool](a, b T) int { return cmp.CompareBool(a, b) }

func CompareString[T ~string](a, b T) int { return cmp.CompareString(a, b) }
func CompareBytes[T ~[]byte](a, b T) int  { return cmp.CompareBytes(a, b) }

func CompareInteger[T Integer](a, b T) int { return cmp.CompareComparable(a, b) }
func CompareFloat[T Float](a, b T) int     { return cmp.CompareComparable(a, b) }
func CompareNumber[T Number](a, b T) int   { return cmp.CompareComparable(a, b) }
func CompareComplex[T Complex](a, b T) int { return cmp.CompareComplex(a, b) }

func CompareComparer[T Comparer[T]](a, b T) int  { return cmp.CompareComparer(a, b) }
func CompareComparable[T Comparable](a, b T) int { return cmp.CompareComparable(a, b) }

// Deprecated: use `CompareComparable` instead.
func CompareOrderable[T Comparable](a, b T) int { return cmp.CompareComparable(a, b) }

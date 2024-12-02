/*
 * Copyright (C) distroy
 */

package cmp

import (
	"bytes"

	"github.com/distroy/ldgo/v2/ldconv"
)

func CompareInterface(a, b interface{}) int {
	aa := reflectValueOf(a)
	bb := reflectValueOf(b)
	return CompareReflect(aa, bb)
}

func CompareBool[T ~bool](a, b T) int {
	switch {
	case a == b:
		return 0
	case bool(a):
		return 1
	default:
		return -1
	}
}

func CompareComplex[T ~complex64 | ~complex128](aa, bb T) int {
	a, b := complex128(aa), complex128(bb)
	if r := CompareOrderable(real(a), real(b)); r != 0 {
		return r
	}
	return CompareOrderable(imag(a), imag(b))
}

func CompareString[T ~string](a, b T) int {
	// aa := *(*string)((unsafe.Pointer)(&a))
	// bb := *(*string)((unsafe.Pointer)(&b))
	aa := string(a)
	bb := string(b)
	return bytes.Compare(ldconv.StrToBytesUnsafe(aa), ldconv.StrToBytesUnsafe(bb))
}

func CompareBytes[T ~[]byte](a, b T) int {
	// aa := *(*[]byte)((unsafe.Pointer)(&a))
	// bb := *(*[]byte)((unsafe.Pointer)(&b))
	aa := []byte(a)
	bb := []byte(b)
	return bytes.Compare(aa, bb)
}

func CompareComparer[T Comparer[T]](a, b T) int { return a.Compare(b) }

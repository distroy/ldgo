/*
 * Copyright (C) distroy
 */

package ldcmp

import (
	"time"

	"github.com/distroy/ldgo/v2/internal/cmp"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Orderable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

type Comparer[T any] interface {
	Compare(b T) int
}

func Compare(a, b interface{}) int          { return cmp.CompareInterface(a, b) }
func CompareInterface(a, b interface{}) int { return cmp.CompareInterface(a, b) }

func CompareBool[T ~bool](a, b T) int { return cmp.CompareBool(a, b) }

func CompareString[T ~string](a, b T) int { return cmp.CompareString(a, b) }
func CompareBytes[T ~[]byte](a, b T) int  { return cmp.CompareBytes(a, b) }

func CompareInteger[T Integer](a, b T) int                  { return CompareOrderable(a, b) }
func CompareFloat[T ~float32 | ~float64](a, b T) int        { return CompareOrderable(a, b) }
func CompareComplex[T ~complex64 | ~complex128](a, b T) int { return cmp.CompareComplex(a, b) }

func CompareComparer[T Comparer[T]](a, b T) int { return cmp.CompareComparer(a, b) }
func CompareOrderable[T Orderable](a, b T) int  { return cmp.CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareByte(a, b byte) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareRune(a, b rune) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareInt(a, b int) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareInt8(a, b int8) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareInt16(a, b int16) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareInt32(a, b int32) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareInt64(a, b int64) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUint(a, b uint) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUint8(a, b uint8) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUint16(a, b uint16) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUint32(a, b uint32) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUint64(a, b uint64) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareUintptr(a, b uintptr) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareFloat32(a, b float32) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareFloat64(a, b float64) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareComplex` instead.
func CompareComplex64(a, b complex64) int { return CompareComplex(a, b) }

// Deprecated: use `CompareComplex` instead.
func CompareComplex128(a, b complex128) int { return CompareComplex(a, b) }

// Deprecated: use `CompareOrderable` instead.
func CompareDuration(a, b time.Duration) int { return CompareOrderable(a, b) }

// Deprecated: use `CompareComparer` instead.
func CompareTime(a, b time.Time) int { return CompareComparer(a, b) }

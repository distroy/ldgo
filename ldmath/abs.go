/*
 * Copyright (C) distroy
 */

package ldmath

type Absable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Abs[T Absable](n T) T {
	if n < 0 {
		n = -n
	}
	return n
}

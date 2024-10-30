/*
 * Copyright (C) distroy
 */

package ldmath

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Max[T Number](n T, args ...T) T {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

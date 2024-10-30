/*
 * Copyright (C) distroy
 */

package ldmath

func Sum[T Number](args ...T) T {
	return sum[T](args...)
}

func sum[R Number, T Number](args ...T) R {
	var sum R = 0
	for _, v := range args {
		sum += R(v)
	}
	return sum
}

func SumInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr](args ...T) int64 {
	return sum[int64](args...)
}

func SumFloat[T float32 | float64](args ...T) float64 {
	return sum[float64](args...)
}

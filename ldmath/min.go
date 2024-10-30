/*
 * Copyright (C) distroy
 */

package ldmath

func Min[T Number](n T, args ...T) T {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

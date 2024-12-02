/*
 * Copyright (C) distroy
 */

package ldslice

import "github.com/distroy/ldgo/v2/internal/cmp"

func Index[S ~[]V, V any](s S, v V) int {
	for i := range s {
		if cmp.CompareInterface(s[i], v) == 0 {
			return i
		}
	}
	return -1
}

func Contains[S ~[]V, V any](s S, v V) bool { return Index(s, v) >= 0 }

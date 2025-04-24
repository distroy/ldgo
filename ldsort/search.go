/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"sort"
)

func internalSearch(n int, f func(i int) bool) int { return sort.Search(n, f) }

func internalIndex(n int, compare func(i int) int) int {
	idx := internalSearch(n, func(i int) bool {
		return compare(i) >= 0
	})
	if idx < n && compare(idx) == 0 {
		return idx
	}
	return -1
}

// Search uses binary search to find and return the smallest index in [0, n)
// at which f(v) == true, or return n if there is no f(v) == true
func Search[T any](a []T, f func(v T) bool) int {
	return internalSearch(len(a), func(i int) bool {
		return f(a[i])
	})
}

// Index uses binary search to find and return the smallest index in [0, n)
// at which compare(v) == 0, or return -1 if there is no compare(v) == 0
//
// To permit binary search, Find requires that compare(i) < 0 for a leading
// prefix of the range, cmp(i) == 0 in the middle, and cmp(i) > 0 for
// the final suffix of the range.
//
// For example:
//
//	Index(a, func (v int) int {
//		if v < target { return -1 }
//		if v > target { return +1 }
//		return 0
//	})
func Index[T any](a []T, compare func(v T) int) int {
	return internalIndex(len(a), func(i int) int {
		return compare(a[i])
	})
}

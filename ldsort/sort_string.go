/*
 * Copyright (C) distroy
 */

package ldsort

import "sort"

type SortSliceString []string

func (s SortSliceString) Len() int           { return len(s) }
func (s SortSliceString) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortSliceString) Less(i, j int) bool { return s[i] < s[j] }

func SortStrings(a []string)          { sort.Sort(SortSliceString(a)) }
func IsSortedStrings(a []string) bool { return sort.IsSorted(SortSliceString(a)) }
func SearchStrings(a []string, x string) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

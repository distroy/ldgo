/*
 * Copyright (C) distroy
 */

package ldsort

type strings = Slice[string]

func SortStrings(a []string)                 { internalSort(strings(a)) }
func UniqStrings(a []string) []string        { return a[:uniq(strings(a))] }
func IsSortedStrings(a []string) bool        { return internalIsSorted(strings(a)) }
func SearchStrings(a []string, x string) int { return templateSearch[string](a, x) }
func IndexStrings(a []string, x string) int  { return templateIndex[string](a, x) }

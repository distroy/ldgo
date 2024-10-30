/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"os"
	"testing"

	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/distroy/ldgo/v3/ldsort"
)

const (
	_count = 100
)

var (
	_numsUnordered []int
	_nums          []int
)

func TestMain(m *testing.M) {
	_numsUnordered = ldrand.Perm(_count)
	_nums = make([]int, 0, len(_numsUnordered))
	_nums = append(_nums, _numsUnordered...)
	ldsort.SortInts(_nums)

	// log.Printf("unordered: %v", _numsUnordered)
	// log.Printf("nums: %v", _nums)

	os.Exit(m.Run())
}

func testNewRBTree() *RBTree[int] {
	rbtree := &RBTree[int]{}
	for _, n := range _numsUnordered {
		rbtree.Insert(n)
	}
	return rbtree
}

func testRBTreeDeleteAll(rbtree *RBTree[int], d int) {
	for it := rbtree.Search(d); it != rbtree.End(); it = rbtree.Search(d) {
		rbtree.Delete(it)
	}
}

func testRBTreeRDeleteAll(rbtree *RBTree[int], d int) {
	for it := rbtree.RSearch(d); it != rbtree.REnd(); it = rbtree.RSearch(d) {
		rbtree.RDelete(it)
	}
}

func testNewMap() *Map[int, int] {
	m := &Map[int, int]{}
	for _, n := range _numsUnordered {
		m.Insert(n, n)
	}
	return m
}

func testMapDeleteAll(m *Map[int, int], d int) {
	for it := m.Search(d); it != m.End(); it = m.Search(d) {
		m.Delete(it)
	}
}

func testMapRDeleteAll(m *Map[int, int], d int) {
	for it := m.RSearch(d); it != m.REnd(); it = m.RSearch(d) {
		m.RDelete(it)
	}
}

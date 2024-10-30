/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo/v3/ldcmp"
	"github.com/distroy/ldgo/v3/ldrand"
)

// cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
// BenchmarkMap_Insert/insert-size:100-12            161498              6897 ns/op
// BenchmarkMap_Insert/insert-size:1000-12            12456             97745 ns/op
// BenchmarkMap_Insert/insert-size:10000-12             775           1396449 ns/op
// BenchmarkMap_Search/search-size:100-12          30981765                41.22 ns/op
// BenchmarkMap_Search/search-size:1000-12         23162068                56.80 ns/op
// BenchmarkMap_Search/search-size:10000-12        12986542                89.42 ns/op

func benchmarkMap_Insert(b *testing.B, size int) {
	b.ResetTimer()
	b.Run(fmt.Sprintf("insert-size:%d", size), func(b *testing.B) {
		// b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				rbtree := &Map[int, int]{
					KeyCompare: func(a, b int) int { return ldcmp.CompareInt(a, b) },
				}
				for i := 0; i < size; i++ {
					rbtree.Insert(ldrand.Int(), 0)
				}
				rbtree.Clear()
			}
		})
	})
	b.StopTimer()
}

func benchmarkMap_Search(b *testing.B, size int) {
	rbtree := &Map[int, int]{
		KeyCompare: func(a, b int) int { return ldcmp.CompareInt(a, b) },
	}
	for i := 0; i < size; i++ {
		rbtree.Insert(ldrand.Int(), 0)
	}
	b.ResetTimer()
	b.Run(fmt.Sprintf("search-size:%d", size), func(b *testing.B) {
		// b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				rbtree.Search(ldrand.Int())
			}
		})
	})
	b.StopTimer()
	rbtree.Clear()
}

func BenchmarkMap_Insert(b *testing.B) {
	benchmarkMap_Insert(b, 100)
	benchmarkMap_Insert(b, 1000)
	benchmarkMap_Insert(b, 10000)
}

func BenchmarkMap_Search(b *testing.B) {
	benchmarkMap_Search(b, 100)
	benchmarkMap_Search(b, 1000)
	benchmarkMap_Search(b, 10000)
}

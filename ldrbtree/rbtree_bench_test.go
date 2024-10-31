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
// BenchmarkRBTree_Insert/insert-size:100-12                 224559              5474 ns/op
// BenchmarkRBTree_Insert/insert-size:1000-12                 19231             60777 ns/op
// BenchmarkRBTree_Insert/insert-size:10000-12                 1258            798311 ns/op
// BenchmarkRBTree_Search/search-size:100-12               43313860                25.84 ns/op
// BenchmarkRBTree_Search/search-size:1000-12              35610832                32.71 ns/op
// BenchmarkRBTree_Search/search-size:10000-12             22626291                45.91 ns/op

func benchmarkRBTree_Insert(b *testing.B, size int) {
	b.ResetTimer()
	b.Run(fmt.Sprintf("insert-size:%d", size), func(b *testing.B) {
		// b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				rbtree := &RBTree[int]{
					Compare: func(a, b int) int { return ldcmp.CompareComparable(a, b) },
				}
				for i := 0; i < size; i++ {
					rbtree.Insert(ldrand.Int())
				}
				rbtree.Clear()
			}
		})
	})
	b.StopTimer()
}

func benchmarkRBTree_Search(b *testing.B, size int) {
	rbtree := &RBTree[int]{
		Compare: func(a, b int) int { return ldcmp.CompareComparable(a, b) },
	}
	for i := 0; i < size; i++ {
		rbtree.Insert(ldrand.Int())
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

func BenchmarkRBTree_Insert(b *testing.B) {
	benchmarkRBTree_Insert(b, 100)
	benchmarkRBTree_Insert(b, 1000)
	benchmarkRBTree_Insert(b, 10000)
}

func BenchmarkRBTree_Search(b *testing.B) {
	benchmarkRBTree_Search(b, 100)
	benchmarkRBTree_Search(b, 1000)
	benchmarkRBTree_Search(b, 10000)
}

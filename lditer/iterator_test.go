/*
 * Copyright (C) distroy
 */

package lditer

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type testPair[K, V any] struct {
	key   K
	value V
}

func testReadSeq[T any](iter iter.Seq[T], yield func(v T) bool) []T {
	res := make([]T, 0, 16)
	for v := range iter {
		if !yield(v) {
			break
		}
		res = append(res, v)
	}
	return res
}

func testReadSeq2[K, V any](iter iter.Seq2[K, V], yield func(k K, v V) bool) []testPair[K, V] {
	res := make([]testPair[K, V], 0, 16)
	for k, v := range iter {
		if !yield(k, v) {
			break
		}
		res = append(res, testPair[K, V]{k, v})
	}
	return res
}

func TestChan(t *testing.T) {
	tests := []struct {
		name  string
		yield func(v int) bool
		slice []int
		want  []int
	}{
		{
			name:  "v > 0",
			yield: func(v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []int{1, 2, 3, 4},
		},
	}

	fnMakeChan := func(s []int) chan int {
		ch := make(chan int, len(s))
		for _, n := range s {
			ch <- n
		}
		close(ch)
		return ch
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			ch := fnMakeChan(tt.slice)
			got := testReadSeq(Chan(ch), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestChan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChan2(t *testing.T) {
	tests := []struct {
		name  string
		yield func(k, v int) bool
		slice []int
		want  []testPair[int, int]
	}{
		{
			name:  "v > 0",
			yield: func(i, v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []testPair[int, int]{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
	}

	fnMakeChan := func(s []int) chan int {
		ch := make(chan int, len(s))
		for _, n := range s {
			ch <- n
		}
		close(ch)
		return ch
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			ch := fnMakeChan(tt.slice)
			got := testReadSeq2(Chan2(ch), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestChan2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name  string
		yield func(v int) bool
		slice []int
		want  []int
	}{
		{
			name:  "v > 0",
			yield: func(v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []int{1, 2, 3, 4},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq(Slice(tt.slice), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice2(t *testing.T) {
	tests := []struct {
		name  string
		yield func(i, v int) bool
		slice []int
		want  []testPair[int, int]
	}{
		{
			name:  "v > 0",
			yield: func(i, v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []testPair[int, int]{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq2(Slice2(tt.slice), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestSlice2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name  string
		yield func(k int, v string) bool
		args  map[int]string
		want  []testPair[int, string]
	}{
		{
			name:  "true",
			yield: func(k int, v string) bool { return true },
			args:  map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:  []testPair[int, string]{{1, "a"}, {2, "b"}, {3, "x"}, {4, "z"}},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq2(Map(tt.args), tt.yield)
			slices.SortFunc(got, func(a, b testPair[int, string]) int { return a.key - b.key })
			slices.SortFunc(tt.want, func(a, b testPair[int, string]) int { return a.key - b.key })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	tests := []struct {
		name  string
		yield func(k int) bool
		args  map[int]string
		want  []int
	}{
		{
			name:  "true",
			yield: func(k int) bool { return true },
			args:  map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:  []int{1, 2, 3, 4},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq(MapKeys(tt.args), tt.yield)
			slices.SortFunc(got, func(a, b int) int { return a - b })
			slices.SortFunc(tt.want, func(a, b int) int { return a - b })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapValues(t *testing.T) {
	tests := []struct {
		name  string
		yield func(v string) bool
		args  map[int]string
		want  []string
	}{
		{
			name:  "true",
			yield: func(v string) bool { return true },
			args:  map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:  []string{"a", "b", "x", "z"},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq(MapValues(tt.args), tt.yield)
			slices.SortFunc(got, func(a, b string) int { return strings.Compare(a, b) })
			slices.SortFunc(tt.want, func(a, b string) int { return strings.Compare(a, b) })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMapValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

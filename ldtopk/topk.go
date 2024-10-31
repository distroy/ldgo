/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"sync"

	"github.com/distroy/ldgo/v3/ldcmp"
)

type LessFunc[T any] func(a, b T) bool

func defaultLess[T any](a, b T) bool {
	return ldcmp.Compare(a, b) < 0
}

func New[T any](k int, less func(a, b T) bool) *Topk[T] {
	p := &Topk[T]{
		Size: k,
		Less: less,
	}
	return p
}

// Topk will keep at most size elements for which less returns true
type Topk[T any] struct {
	Size   int         // K
	Less   LessFunc[T] // less func
	Locker sync.Locker // locker
	data   []T         //
}

func (p *Topk[T]) Data() []T { return p.data }

func (p *Topk[T]) Add(d T) bool {
	locker := p.Locker
	if locker == nil {
		p.Locker = nullLocker{}
		locker = nullLocker{}
	}

	locker.Lock()
	defer locker.Unlock()

	k := p.Size
	if k <= 0 {
		return false
	}

	p.init()
	if pos := len(p.data); pos < k {
		topkAppendTail(&p.data, p.Less, d)
		return true
	}

	if !p.Less(d, p.data[0]) {
		return false
	}
	p.data[0] = d

	topkFixupHead(&p.data, p.Less)
	return true
}

func (p *Topk[T]) init() {
	if p.Less == nil {
		p.Less = defaultLess[T]
	}
	if p.data == nil && p.Size > 0 {
		p.data = make([]T, 0, p.Size)
	}
}

func topkAppendTail[T any](heap *[]T, less LessFunc[T], d T) {
	pos := len(*heap)

	*heap = append(*heap, d)
	for parent := 0; pos > 0; pos = parent {
		parent = (pos - 1) / 2
		if !less((*heap)[parent], (*heap)[pos]) {
			break
		}
		(*heap)[parent], (*heap)[pos] = (*heap)[pos], (*heap)[parent]
	}
}

func topkFixupHead[T any](heap *[]T, less LessFunc[T]) {
	size := len(*heap)
	for pos := 0; ; {
		lChild := (pos * 2) + 1
		rChild := (pos * 2) + 2
		if lChild >= size {
			break
		}

		child := lChild
		if rChild < size && less((*heap)[lChild], (*heap)[rChild]) {
			child = rChild
		}

		if !less((*heap)[pos], (*heap)[child]) {
			break
		}

		(*heap)[pos], (*heap)[child] = (*heap)[child], (*heap)[pos]
		pos = child
	}
}

type nullLocker struct{}

func (_ nullLocker) Lock()   {}
func (_ nullLocker) Unlock() {}

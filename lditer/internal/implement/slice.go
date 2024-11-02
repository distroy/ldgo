/*
 * Copyright (C) distroy
 */

package implement

import (
	"github.com/distroy/ldgo/v3/lditer"
	"github.com/distroy/ldgo/v3/ldslice"
)

var (
	_ lditer.Iterator[int] = lditer.Iter[int](ldslice.Begin[int]([]int{}))
	_ lditer.Iterator[int] = lditer.Iter[int](ldslice.End[int]([]int{}))

	_ lditer.Iterator[int] = lditer.Iter[int](ldslice.RBegin[int]([]int{}))
	_ lditer.Iterator[int] = lditer.Iter[int](ldslice.REnd[int]([]int{}))

	_ lditer.Range[int] = (*ldslice.Range[int])(nil)
	_ lditer.Range[int] = (*ldslice.ReverseRange[int])(nil)
)

var (
	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldslice.Begin[int]([]int{}))
	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldslice.End[int]([]int{}))

	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldslice.RBegin[int]([]int{}))
	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldslice.REnd[int]([]int{}))

	_ lditer.ConstRange[int] = (*ldslice.Range[int])(nil)
	_ lditer.ConstRange[int] = (*ldslice.ReverseRange[int])(nil)
)

/*
 * Copyright (C) distroy
 */

package implement

import (
	"github.com/distroy/ldgo/v3/lditer"
	"github.com/distroy/ldgo/v3/ldrbtree"
)

var (
	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldrbtree.RBTreeIterator[int]{})
	_ lditer.ConstIterator[int] = lditer.ConstIter[int](ldrbtree.RBTreeReverseIterator[int]{})

	_ lditer.ConstIterator[ldrbtree.MapNode[string, int]] = lditer.ConstIter[ldrbtree.MapNode[string, int]](ldrbtree.MapIterator[string, int]{})
	_ lditer.ConstIterator[ldrbtree.MapNode[string, int]] = lditer.ConstIter[ldrbtree.MapNode[string, int]](ldrbtree.MapReverseIterator[string, int]{})

	_ lditer.ConstRange[int] = (*ldrbtree.RBTreeRange[int])(nil)
	_ lditer.ConstRange[int] = (*ldrbtree.RBTreeReverseRange[int])(nil)

	_ lditer.ConstRange[ldrbtree.MapNode[string, int]] = (*ldrbtree.MapRange[string, int])(nil)
	_ lditer.ConstRange[ldrbtree.MapNode[string, int]] = (*ldrbtree.MapReverseRange[string, int])(nil)
)

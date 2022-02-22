/*
 * Copyright (C) distroy
 */

package ldrbtree

type CompareFunc = func(a, b interface{}) int

// RBTree is red-black tree
type RBTree struct {
	Compare  CompareFunc
	root     *rbtreeNode
	sentinel *rbtreeNode
	count    int
}

func (rbt *RBTree) Len() int {
	return rbt.count
}

func (rbt *RBTree) Insert(d interface{}) RBTreeIterator {
	rbt.init()

	root := &rbt.root

	node := getRBTreeNode(rbt.sentinel)
	node.Data = d
	it := RBTreeIterator(rbtreeIterator{
		tree: rbt,
		node: node,
	})

	if *root == rbt.sentinel {
		*root = node
		rbt.count++
		return it
	}

	rbt.insertNode(node)
	rbt.insertFixup(node)
	rbt.count++
	return it
}

func (rbt *RBTree) Clear() {
	rbt.init()

	node := rbt.root
	sentinel := rbt.sentinel

	if node == sentinel {
		return
	}

	buffer := make([]*rbtreeNode, 0, (rbt.count+1)/2)
	buffer = append(buffer, node)
	rbt.root = sentinel
	rbt.count = 0

	for len(buffer) > 0 {
		pos := len(buffer) - 1
		node := buffer[pos]
		buffer = buffer[:pos]

		if node.Left != sentinel {
			buffer = append(buffer, node.Left)
		}
		if node.Right != sentinel {
			buffer = append(buffer, node.Right)
		}
		putRBTreeNode(node)
	}
}

func (rbt *RBTree) Delete(it RBTreeIterator) RBTreeIterator {
	rbt.init()

	if it.tree != rbt {
		panic("the iterator does not belong to the red-black tree, can not delete")
	}

	sentinel := rbt.sentinel
	node := it.node
	if node == sentinel {
		panic("the iterator is already at the end of the red-black tree, can not delete")
	}

	it = it.Next()
	rbt.deleteNode(node)
	rbt.count--

	putRBTreeNode(node)
	return it
}

func (rbt *RBTree) Search(d interface{}) RBTreeIterator {
	rbt.init()

	sentinel := rbt.sentinel

	node := rbt.lowerBound(d, forward(rbt))
	if node != sentinel && rbt.Compare(d, node.Data) != 0 {
		node = sentinel
	}
	return RBTreeIterator(rbtreeIterator{
		tree: rbt,
		node: node,
	})
}

func (rbt *RBTree) Range() *RBTreeRange {
	rbt.init()

	return &RBTreeRange{
		begin: rbt.Begin(),
		end:   rbt.End(),
	}
}

func (rbt *RBTree) Begin() RBTreeIterator {
	rbt.init()
	return RBTreeIterator(rbt.beginIterator(forward(rbt)))
}

func (rbt *RBTree) End() RBTreeIterator {
	rbt.init()
	return RBTreeIterator(rbt.endIterator(forward(rbt)))
}

// RDelete is reverse delete
func (rbt *RBTree) RDelete(it RBTreeReverseIterator) RBTreeReverseIterator {
	rbt.init()

	if it.tree != rbt {
		panic("the iterator does not belong to the red-black tree, can not delete")
	}

	sentinel := rbt.sentinel
	node := it.node
	if node == sentinel {
		panic("the iterator is already at the end of the red-black tree, can not delete")
	}

	it = it.Next()
	rbt.deleteNode(node)
	rbt.count--

	putRBTreeNode(node)
	return it
}

// RSearch is reverse search
func (rbt *RBTree) RSearch(d interface{}) RBTreeReverseIterator {
	rbt.init()

	sentinel := rbt.sentinel

	node := rbt.lowerBound(d, reverse(rbt))
	if node != sentinel && rbt.Compare(d, node.Data) != 0 {
		node = sentinel
	}
	return RBTreeReverseIterator(rbtreeIterator{
		tree: rbt,
		node: node,
	})
}

// RRange is reverse range
func (rbt *RBTree) RRange() *RBTreeReverseRange {
	rbt.init()

	return &RBTreeReverseRange{
		begin: rbt.RBegin(),
		end:   rbt.REnd(),
	}
}

// RBegin is reverse begin
func (rbt *RBTree) RBegin() RBTreeReverseIterator {
	rbt.init()
	return RBTreeReverseIterator(rbt.beginIterator(reverse(rbt)))
}

// REnd is reverse end
func (rbt *RBTree) REnd() RBTreeReverseIterator {
	rbt.init()
	return RBTreeReverseIterator(rbt.endIterator(reverse(rbt)))
}

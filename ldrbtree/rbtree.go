/*
 * Copyright (C) distroy
 */

package ldrbtree

type CompareFunc[T any] func(a, b T) int

// RBTree is red-black tree
type RBTree[T any] struct {
	Compare  CompareFunc[T]
	root     *rbtreeNode[T]
	sentinel *rbtreeNode[T]
	count    int
}

func (rbt *RBTree[T]) Len() int {
	return rbt.count
}

func (rbt *RBTree[T]) Insert(d T) RBTreeIterator[T] {
	rbt.init()

	root := &rbt.root

	node := getRBTreeNode(rbt.sentinel)
	node.Data = d
	it := RBTreeIterator[T]{
		tree: rbt,
		node: node,
	}

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

func (rbt *RBTree[T]) InsertOrSearch(d T) RBTreeIterator[T] {
	rbt.init()

	root := &rbt.root

	node := getRBTreeNode[T](rbt.sentinel)
	node.Data = d

	if *root == rbt.sentinel {
		*root = node
		rbt.count++
		return RBTreeIterator[T]{
			tree: rbt,
			node: node,
		}
	}

	temp := rbt.insertOrSearchNode(node)
	if temp != node {
		// d already exists
		putRBTreeNode[T](node)

	} else {
		// d inserted just now
		rbt.insertFixup(node)
		rbt.count++
	}

	return RBTreeIterator[T]{
		tree: rbt,
		node: temp,
	}
}

func (rbt *RBTree[T]) InsertOrAssign(d T) RBTreeIterator[T] {
	rbt.init()

	root := &rbt.root

	node := getRBTreeNode[T](rbt.sentinel)
	node.Data = d

	if *root == rbt.sentinel {
		*root = node
		rbt.count++
		return RBTreeIterator[T]{
			tree: rbt,
			node: node,
		}
	}

	temp := rbt.insertOrSearchNode(node)
	if temp != node {
		// d already exists
		temp.Data = d
		putRBTreeNode[T](node)

	} else {
		// d inserted just now
		rbt.insertFixup(node)
		rbt.count++
	}

	return RBTreeIterator[T]{
		tree: rbt,
		node: temp,
	}
}

func (rbt *RBTree[T]) Clear() {
	rbt.init()

	node := rbt.root
	sentinel := rbt.sentinel

	rbt.root = sentinel
	rbt.count = 0

	for node != sentinel {
		if node.Left != sentinel && node.Right != sentinel {
			node = node.Left
			continue
		}

		// node.Left == sentinel || node.Right == sentinel
		parent := node.Parent
		child := node.Left
		if node.Left == sentinel {
			child = node.Right
		}

		if node == parent.Left {
			parent.Left = child
		} else if node == parent.Right {
			parent.Right = child
		}

		if child != sentinel {
			child.Parent = parent
		}

		// ldlog.Default().Info("*** clear", zap.Any("data", node.Data))
		putRBTreeNode[T](node)

		if child != sentinel {
			node = child
		} else {
			node = parent
		}
	}
}

func (rbt *RBTree[T]) Delete(it RBTreeIterator[T]) RBTreeIterator[T] {
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

	putRBTreeNode[T](node)
	return it
}

func (rbt *RBTree[T]) Count(d T) int {
	rbt.init()

	sentinel := rbt.sentinel

	count := 0
	node := rbtreeLowerBound(d, forward(rbt))
	for node != sentinel && rbt.Compare(d, node.Data) == 0 {
		count++
		node = node.next(forward(rbt))
	}
	return count
}

// Search returns the first node == d
func (rbt *RBTree[T]) Search(d T) RBTreeIterator[T] {
	rbt.init()

	sentinel := rbt.sentinel

	node := rbtreeLowerBound(d, forward(rbt))
	if node != sentinel && rbt.Compare(d, node.Data) != 0 {
		node = sentinel
	}
	return RBTreeIterator[T]{
		tree: rbt,
		node: node,
	}
}

// SearchRange returns the node range [first node == d, first node > d)
func (rbt *RBTree[T]) SearchRange(d T) *RBTreeRange[T] {
	rbt.init()

	sentinel := rbt.sentinel

	begin := rbtreeLowerBound(d, forward(rbt))
	if begin != sentinel && rbt.Compare(d, begin.Data) != 0 {
		return &RBTreeRange[T]{
			Begin: RBTreeIterator[T]{tree: rbt, node: sentinel},
			End:   RBTreeIterator[T]{tree: rbt, node: sentinel},
		}
	}

	end := rbtreeUpperBound(d, forward(rbt))
	return &RBTreeRange[T]{
		Begin: RBTreeIterator[T]{tree: rbt, node: begin},
		End:   RBTreeIterator[T]{tree: rbt, node: end},
	}
}

// LowerBound returns the first node >= d
func (rbt *RBTree[T]) LowerBound(d T) RBTreeIterator[T] {
	rbt.init()
	return RBTreeIterator[T]{
		tree: rbt,
		node: rbtreeLowerBound(d, forward(rbt)),
	}
}

// UpperBound returns the first node > d
func (rbt *RBTree[T]) UpperBound(d T) RBTreeIterator[T] {
	rbt.init()
	return RBTreeIterator[T]{
		tree: rbt,
		node: rbtreeUpperBound(d, forward(rbt)),
	}
}

func (rbt *RBTree[T]) Range() *RBTreeRange[T] {
	rbt.init()

	return &RBTreeRange[T]{
		Begin: rbt.Begin(),
		End:   rbt.End(),
	}
}

func (rbt *RBTree[T]) Begin() RBTreeIterator[T] {
	rbt.init()
	return RBTreeIterator[T](rbtreeBeginIterator[T](forward(rbt)))
}

func (rbt *RBTree[T]) End() RBTreeIterator[T] {
	rbt.init()
	return RBTreeIterator[T](rbtreeEndIterator(forward(rbt)))
}

// RDelete is reverse delete
func (rbt *RBTree[T]) RDelete(it RBTreeReverseIterator[T]) RBTreeReverseIterator[T] {
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

	putRBTreeNode[T](node)
	return it
}

// RSearch is reverse search
// RSearch returns the last node == d
func (rbt *RBTree[T]) RSearch(d T) RBTreeReverseIterator[T] {
	rbt.init()

	sentinel := rbt.sentinel

	node := rbtreeLowerBound(d, reverse(rbt))
	if node != sentinel && rbt.Compare(d, node.Data) != 0 {
		node = sentinel
	}
	return RBTreeReverseIterator[T]{
		tree: rbt,
		node: node,
	}
}

// RSearchRange is reverse search range
// RSearchRange returns the node range [last node == d, last node < d)
func (rbt *RBTree[T]) RSearchRange(d T) *RBTreeReverseRange[T] {
	rbt.init()

	sentinel := rbt.sentinel

	begin := rbtreeLowerBound(d, reverse(rbt))
	if begin != sentinel && rbt.Compare(d, begin.Data) != 0 {
		return &RBTreeReverseRange[T]{
			Begin: RBTreeReverseIterator[T]{tree: rbt, node: sentinel},
			End:   RBTreeReverseIterator[T]{tree: rbt, node: sentinel},
		}
	}

	end := rbtreeUpperBound(d, reverse(rbt))
	return &RBTreeReverseRange[T]{
		Begin: RBTreeReverseIterator[T]{tree: rbt, node: begin},
		End:   RBTreeReverseIterator[T]{tree: rbt, node: end},
	}
}

// RLowerBound is reverse lower bound
// RLowerBound returns the last node <= d
func (rbt *RBTree[T]) RLowerBound(d T) RBTreeReverseIterator[T] {
	rbt.init()
	return RBTreeReverseIterator[T]{
		tree: rbt,
		node: rbtreeLowerBound(d, reverse(rbt)),
	}
}

// RUpperBound is reverse upper bound
// RUpperBound returns the last node < d
func (rbt *RBTree[T]) RUpperBound(d T) RBTreeReverseIterator[T] {
	rbt.init()
	return RBTreeReverseIterator[T]{
		tree: rbt,
		node: rbtreeUpperBound(d, reverse(rbt)),
	}
}

// RRange is reverse range
func (rbt *RBTree[T]) RRange() *RBTreeReverseRange[T] {
	rbt.init()

	return &RBTreeReverseRange[T]{
		Begin: rbt.RBegin(),
		End:   rbt.REnd(),
	}
}

// RBegin is reverse begin
func (rbt *RBTree[T]) RBegin() RBTreeReverseIterator[T] {
	rbt.init()
	return RBTreeReverseIterator[T](rbtreeBeginIterator[T](reverse(rbt)))
}

// REnd is reverse end
func (rbt *RBTree[T]) REnd() RBTreeReverseIterator[T] {
	rbt.init()
	return RBTreeReverseIterator[T](rbtreeEndIterator(reverse(rbt)))
}

/*
 * Copyright (C) distroy
 */

package ldrbtree

func (rbt *RBTree[T]) insertNode(node *rbtreeNode[T]) {
	sentinel := rbt.sentinel

	p := &rbt.root
	last := *p
	for {
		r := rbt.Compare(node.Data, (*p).Data)
		if r < 0 {
			p = &(*p).Left
		} else {
			p = &(*p).Right
		}

		if *p == sentinel {
			break
		}
		last = *p
	}

	*p = node
	node.Parent = last
	node.Color = colorRed

	// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
}

func (rbt *RBTree[T]) insertOrSearchNode(node *rbtreeNode[T]) *rbtreeNode[T] {
	sentinel := rbt.sentinel

	p := &rbt.root
	last := *p
	for {
		r := rbt.Compare(node.Data, (*p).Data)
		if r == 0 {
			return *p
		} else if r < 0 {
			p = &(*p).Left
		} else {
			p = &(*p).Right
		}

		if *p == sentinel {
			break
		}
		last = *p
	}

	*p = node
	node.Parent = last
	node.Color = colorRed

	// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
	return node
}

func (rbt *RBTree[T]) insertFixup(node *rbtreeNode[T]) {
	root := &rbt.root

	for node != *root && node.Parent.Color == colorRed {
		if node.Parent == node.Parent.Parent.Left {
			rbt.insertFixupLeftParent(&node)
		} else {
			rbt.insertFixupRightParent(&node)
		}
	}

	(*root).Color = colorBlack
}

func (rbt *RBTree[T]) insertFixupLeftParent(node **rbtreeNode[T]) {
	uncle := (*node).Parent.Parent.Right
	if uncle.Color == colorRed {
		// case 1: uncle is red
		// (color flips)
		//       G            g
		//      / \          / \
		//     p   u  -->   P   U
		//    /            /
		//   n            n
		(*node).Parent.Color = colorBlack
		uncle.Color = colorBlack
		(*node).Parent.Parent.Color = colorRed
		(*node) = (*node).Parent.Parent
		// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
		return
	}

	if (*node) == (*node).Parent.Right {
		// case 2: uncle is black and node is the parent's right child
		// (left rotate at parent)
		//     G             G
		//    / \           / \
		//   p   U  -->    n   U
		//    \           /
		//     n         p
		(*node) = (*node).Parent
		rbt.rotateLeft(*node)
		// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
	}

	// case 3: uncle is black and node is the parent's left child
	// (right rotate at gparent).
	//       G           P
	//      / \         / \
	//     p   U  -->  n   g
	//    /                 \
	//   n                   U
	(*node).Parent.Color = colorBlack
	(*node).Parent.Parent.Color = colorRed
	rbt.rotateRight((*node).Parent.Parent)
	// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
}

func (rbt *RBTree[T]) insertFixupRightParent(node **rbtreeNode[T]) {
	uncle := (*node).Parent.Parent.Left
	if uncle.Color == colorRed {
		// case 1
		(*node).Parent.Color = colorBlack
		uncle.Color = colorBlack
		(*node).Parent.Parent.Color = colorRed
		(*node) = (*node).Parent.Parent
		// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
		return
	}

	// case 2
	if (*node) == (*node).Parent.Left {
		(*node) = (*node).Parent
		rbt.rotateRight(*node)
		// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
	}

	// case 3
	(*node).Parent.Color = colorBlack
	(*node).Parent.Parent.Color = colorRed
	rbt.rotateLeft((*node).Parent.Parent)
	// ldlog.Default().Info("==== debug", ldlog.Reflect("tree", rbt.toMap()))
}

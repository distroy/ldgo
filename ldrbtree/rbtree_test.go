/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo/v3/ldcmp"
	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/smartystreets/goconvey/convey"
)

func TestRBTree_Insert(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("insert", func(c convey.C) {
			rbtree := &RBTree[int]{}
			for _, n := range _numsUnordered {
				rbtree.Insert(n)
				c.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				c.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)
			}

			c.Convey("len", func(c convey.C) {
				c.So(rbtree.Len(), convey.ShouldEqual, len(_nums))
			})

			for _, n := range _numsUnordered {
				rbtree.Insert(n)
				c.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				c.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)
			}
			c.Convey("twice len", func(c convey.C) {
				c.So(rbtree.Len(), convey.ShouldEqual, len(_nums)*2)
			})

			c.Convey("check parent", func(c convey.C) {
				c.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
			})

			c.Convey("check color", func(c convey.C) {
				c.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)
			})

			c.Convey("clear", func(c convey.C) {
				rbtree.Clear()
				c.So(rbtree.Len(), convey.ShouldEqual, 0)
				c.So(rbtree.Begin(), convey.ShouldResemble, rbtree.End())
			})
		})

		c.Convey("insert or search", func(c convey.C) {
			rbtree := testNewRBTree()
			c.So(rbtree.Len(), convey.ShouldEqual, len(_nums))

			c.Convey("insert when exists", func(c convey.C) {
				for _, n := range _numsUnordered {
					it := rbtree.InsertOrSearch(n)
					c.So(it.Get(), convey.ShouldEqual, n)
					c.So(rbtree.Len(), convey.ShouldEqual, len(_nums))
				}
				it := rbtree.Begin()
				for _, n := range _nums {
					c.So(it.Get(), convey.ShouldEqual, n)
					it = it.Next()
				}
				c.So(it, convey.ShouldResemble, rbtree.End())
			})

			c.Convey("insert when not exists", func(c convey.C) {
				it := rbtree.InsertOrSearch(-100)
				c.So(it.Get(), convey.ShouldEqual, -100)
				c.So(rbtree.Len(), convey.ShouldEqual, len(_nums)+1)
			})

			c.Convey("insert when empty", func(c convey.C) {
				rbtree.Clear()
				c.So(rbtree.Len(), convey.ShouldEqual, 0)

				it := rbtree.InsertOrSearch(100)
				c.So(it.Get(), convey.ShouldEqual, 100)
				c.So(rbtree.Len(), convey.ShouldEqual, 1)
			})
		})

		c.Convey("insert or assign", func(c convey.C) {
			rbtree := &RBTree[[2]int]{
				Compare: func(a, b [2]int) int {
					aa, bb := a, b
					return ldcmp.CompareComparable(aa[0], bb[0])
				},
			}
			c.So(rbtree.Len(), convey.ShouldEqual, 0)

			rbtree.InsertOrAssign([2]int{100, 0})
			c.So(rbtree.Len(), convey.ShouldEqual, 1)
			c.So(rbtree.Begin().Get(), convey.ShouldResemble, [2]int{100, 0})
			c.So(rbtree.Begin().Next(), convey.ShouldResemble, rbtree.End())

			rbtree.InsertOrAssign([2]int{100, 2})
			c.So(rbtree.Len(), convey.ShouldEqual, 1)
			c.So(rbtree.Begin().Get(), convey.ShouldResemble, [2]int{100, 2})
			c.So(rbtree.Begin().Next(), convey.ShouldResemble, rbtree.End())

			rbtree.InsertOrAssign([2]int{200, 4})
			c.So(rbtree.Len(), convey.ShouldEqual, 2)
			c.So(rbtree.Begin().Get(), convey.ShouldResemble, [2]int{100, 2})
			c.So(rbtree.RBegin().Get(), convey.ShouldResemble, [2]int{200, 4})
		})
	})
}

func TestRBTree_DuplicateData(t *testing.T) {
	fnDeleteAll := testRBTreeDeleteAll

	convey.Convey(t.Name(), t, func(c convey.C) {
		const retry = 20
		c.Convey("search", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.Search(d)

				c.So(it.Get(), convey.ShouldEqual, d)
				c.So(it.Next().Get(), convey.ShouldEqual, d)
				c.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})

		c.Convey("lower bound", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.LowerBound(d)

				c.So(it.Get(), convey.ShouldEqual, d)
				c.So(it.Next().Get(), convey.ShouldEqual, d)
				c.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})

		c.Convey("upper bound", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)

				rbtree.Insert(d - 1)
				rbtree.Insert(d)
				rbtree.Insert(d + 1)

				it := rbtree.UpperBound(d)

				c.So(it.Get(), convey.ShouldBeGreaterThan, d)
				c.So(it.Prev().Get(), convey.ShouldBeLessThanOrEqualTo, d)
				// convey.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})
	})
}

func TestRBTree_Count(t *testing.T) {
	const retry = 20
	convey.Convey(t.Name(), t, func(c convey.C) {
		rbtree := testNewRBTree()
		for i := 0; i < retry; i++ {
			d := ldrand.Intn(_count)

			count0 := 0
			for it := rbtree.SearchRange(d); it.HasNext(); it.Next() {
				count0++
			}

			c.So(rbtree.Count(d), convey.ShouldEqual, count0)

			rbtree.Insert(d)
			rbtree.Insert(d)
			rbtree.Insert(d)

			c.So(rbtree.Count(d), convey.ShouldEqual, count0+3)
		}

		c.So(rbtree.Count(-100), convey.ShouldEqual, 0)
	})
}

func TestRBTree_DuplicateDataReverse(t *testing.T) {
	fnDeleteAll := testRBTreeRDeleteAll
	const retry = 20

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("search", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.RSearch(d)

				c.So(it.Get(), convey.ShouldEqual, d)
				c.So(it.Next().Get(), convey.ShouldEqual, d)
				c.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})

		c.Convey("lower bound", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)
				fnDeleteAll(rbtree, d)

				rbtree.Insert(d)
				rbtree.Insert(d)
				rbtree.Insert(d)

				it := rbtree.RLowerBound(d)

				c.So(it.Get(), convey.ShouldEqual, d)
				c.So(it.Next().Get(), convey.ShouldEqual, d)
				c.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})

		c.Convey("upper bound", func(c convey.C) {
			for i := 0; i < retry; i++ {
				rbtree := testNewRBTree()

				d := ldrand.Intn(_count)
				rbtree.Insert(d - 1)
				rbtree.Insert(d)
				rbtree.Insert(d + 1)

				it := rbtree.RUpperBound(d)

				c.So(it.Get(), convey.ShouldBeLessThan, d)
				c.So(it.Prev().Get(), convey.ShouldBeGreaterThanOrEqualTo, d)
				// convey.So(it.Next().Next().Get(), convey.ShouldEqual, d)
			}
		})
	})
}

func (n *rbtreeNode[T]) checkParent(sentinel *rbtreeNode[T]) bool {
	if n == sentinel {
		return true
	}
	if n.Parent != sentinel {
		if n.Parent.Left != n && n.Parent.Right != n {
			// ldlog.Default().Info("check parent fail", ldlog.Any("node", n.Data), ldlog.Any("parent", n.Parent.Data), ldlog.Uintptr("sentinel", uintptr(unsafe.Pointer(sentinel))))
			return false
		}
	}

	if n.Left != sentinel {
		if n.Left.Parent != n {
			// ldlog.Default().Info("check parent fail", ldlog.Any("node", n.Data), ldlog.Any("left", n.Left.Data))
			return false
		} else if !n.Left.checkParent(sentinel) {
			return false
		}
	}

	if n.Right != sentinel {
		if n.Right.Parent != n {
			// ldlog.Default().Info("check parent fail", ldlog.Any("node", n.Data), ldlog.Any("right", n.Right.Data))
			return false
		} else if !n.Right.checkParent(sentinel) {
			return false
		}
	}

	return true
}

func (n *rbtreeNode[T]) checkColor(sentinel *rbtreeNode[T]) bool {
	blacks := n.getBlackCount(sentinel)
	return n.walkColor(sentinel, blacks, 0)
}

func (n *rbtreeNode[T]) walkColor(sentinel *rbtreeNode[T], wantBlacks, currentBlacks int) bool {
	if n.Color == colorRed {
		if n.Left == nil || n.Right == nil || n.Left.Color != colorBlack || n.Right.Color != colorBlack {
			// ldlog.Default().Info("the children of red node is not black", ldlog.Reflect("node", n.toMap(sentinel)))
			return false
		}
	}

	if n.Color == colorBlack {
		currentBlacks++
	}

	if n.Left == nil {
		return currentBlacks == wantBlacks
	}

	if !n.Left.walkColor(sentinel, wantBlacks, currentBlacks) {
		return false
	}
	if !n.Right.walkColor(sentinel, wantBlacks, currentBlacks) {
		return false
	}

	return true
}

func (n *rbtreeNode[T]) getBlackCount(sentinel *rbtreeNode[T]) int {
	blacks := 0
	for node := n; node != nil; node = node.Left {
		if node.Color == colorBlack {
			blacks++
		}
	}
	return blacks
}

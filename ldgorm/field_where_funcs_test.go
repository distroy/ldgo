/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo/v3/ldptr"
	"github.com/smartystreets/goconvey/convey"
)

func TestBetween(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Between(nil, nil)", func() {
			cond := Between(nil, nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Between(int, nil)", func() {
			cond := Between(0, nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Between(string, nil)", func() {
			cond := Between("abc", nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Between(*int, nil)", func() {
			cond := Between(ldptr.New[int32](0), nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Between(*int(nil), nil)", func() {
			cond := Between((*int)(nil), nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Between(*string, nil)", func() {
			cond := Between(ldptr.New("abc"), nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Between(*string(nil), nil)", func() {
			cond := Between((*string)(nil), nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Between(nil, int)", func() {
			cond := Between(nil, 0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Between(nil, string)", func() {
			cond := Between(nil, "abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Between(nil, *int)", func() {
			cond := Between(nil, ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Between(nil, *int(nil))", func() {
			cond := Between(nil, (*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Between(nil, *string)", func() {
			cond := Between(nil, ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Between(nil, *string(nil))", func() {
			cond := Between(nil, (*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Between(int, int)", func() {
			cond := Between(0, 10)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " BETWEEN ? AND ?"), Args: []interface{}{0, 10},
			})
		})
		convey.Convey("Between(int, int) && min == max", func() {
			cond := Between(3, 3)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{3},
			})
		})
		convey.Convey("Between(*int, *int)", func() {
			cond := Between(ldptr.New[int32](0), ldptr.New[int32](10))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " BETWEEN ? AND ?"), Args: []interface{}{int32(0), int32(10)},
			})
		})
		convey.Convey("Between(*int, *int) && min == max", func() {
			cond := Between(ldptr.New[int32](3), ldptr.New[int32](3))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{int32(3)},
			})
		})
		convey.Convey("Between(*int(nil), *int)", func() {
			cond := Between((*int)(nil), ldptr.New[int32](10))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{int32(10)},
			})
		})
		convey.Convey("Between(*int, *int(nil))", func() {
			cond := Between(ldptr.New[int32](0), (*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Between(*int(nil), *int(nil))", func() {
			cond := Between((*int)(nil), (*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

	})
}

func TestGt(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Gt(nil)", func() {
			cond := Gt(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Gt(int)", func() {
			cond := Gt(0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " > ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Gt(string)", func() {
			cond := Gt("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " > ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Gt(*int)", func() {
			cond := Gt(ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " > ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Gt(*int(nil))", func() {
			cond := Gt((*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Gt(*string)", func() {
			cond := Gt(ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " > ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Gt(*string(nil))", func() {
			cond := Gt((*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
	})
}

func TestLt(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Lt(nil)", func() {
			cond := Lt(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Lt(int)", func() {
			cond := Lt(0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " < ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Lt(string)", func() {
			cond := Lt("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " < ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Lt(*int)", func() {
			cond := Lt(ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " < ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Lt(*int(nil))", func() {
			cond := Lt((*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Lt(*string)", func() {
			cond := Lt(ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " < ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Lt(*string(nil))", func() {
			cond := Lt((*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
	})
}

func TestGte(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Gte(nil)", func() {
			cond := Gte(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Gte(int)", func() {
			cond := Gte(0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Gte(string)", func() {
			cond := Gte("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Gte(*int)", func() {
			cond := Gte(ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Gte(*int(nil))", func() {
			cond := Gte((*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Gte(*string)", func() {
			cond := Gte(ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " >= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Gte(*string(nil))", func() {
			cond := Gte((*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
	})
}

func TestLte(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Lte(nil)", func() {
			cond := Lte(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Lte(int)", func() {
			cond := Lte(0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Lte(string)", func() {
			cond := Lte("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Lte(*int)", func() {
			cond := Lte(ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Lte(*int(nil))", func() {
			cond := Lte((*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Lte(*string)", func() {
			cond := Lte(ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " <= ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Lte(*string(nil))", func() {
			cond := Lte((*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
	})
}

func TestEqual(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Equal(nil)", func() {
			cond := Equal(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})

		convey.Convey("Equal(int)", func() {
			cond := Equal(0)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{0},
			})
		})
		convey.Convey("Equal(string)", func() {
			cond := Equal("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Equal(*int)", func() {
			cond := Equal(ldptr.New[int32](0))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{int32(0)},
			})
		})
		convey.Convey("Equal(*int(nil))", func() {
			cond := Equal((*int)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("Equal(*string)", func() {
			cond := Equal(ldptr.New("abc"))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " = ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Equal(*string(nil))", func() {
			cond := Equal((*string)(nil))
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
	})
}

func TestIn(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("In(nil)", func() {
			cond := In(nil)
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldBeZeroValue)
		})
		convey.Convey("In(int)", func() {
			convey.So(func() { In(0) }, convey.ShouldPanic)
		})
		convey.Convey("In([3]int)", func() {
			cond := In([3]int{310, 320, 330})
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " IN (?)"), Args: []interface{}{[3]int{310, 320, 330}},
			})
		})
		convey.Convey("In([]int)", func() {
			cond := In([]int{310, 320, 330})
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " IN (?)"), Args: []interface{}{[]int{310, 320, 330}},
			})
		})
	})
}

func TestLike(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		const field = "field"

		convey.Convey("Like(abc)", func() {
			cond := Like("abc")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " LIKE ?"), Args: []interface{}{"abc"},
			})
		})
		convey.Convey("Like(%abc%)", func() {
			cond := Like("%abc%")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf("%s%s", field, " LIKE ?"), Args: []interface{}{"%abc%"},
			})
		})
	})
}

func TestExpr(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		const field = "field"

		c.Convey(`right({{column}}, 1) = "%"`, func(c convey.C) {
			cond := Expr(`right({{column}}, ?) = ?`, 1, "%")
			res := cond.buildWhere(field)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: fmt.Sprintf(`right(%s, ?) = ?`, field), Args: []interface{}{1, "%"},
			})
		})
	})
}

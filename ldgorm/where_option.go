/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/distroy/ldgo/v3/ldconv"
)

type WhereOption interface {
	Option

	And(where interface{}) WhereOption
	Or(where interface{}) WhereOption

	buildWhere(db *GormDb) whereResult
}

func Where(where interface{}) WhereOption {
	if where == nil {
		panic("the where type must not be nil")
	}

	switch w := where.(type) {
	case *whereOptionTree:
		if len(w.Wheres) == 1 {
			return w
		}
		return &whereOptionTree{
			Wheres: []whereOptionTreeNode{
				{Where: w},
			},
		}

	case *whereOption:
		return w.toTree()
	}

	val := reflect.ValueOf(where)
	w := getWhereReflect(val.Type())

	return &whereOption{
		Iface: where,
		value: val,
		where: w,
	}
}

type whereOption struct {
	Iface interface{}   `json:"where"`
	value reflect.Value `json:"-"`
	where *whereReflect `json:"-"`
}

func (w *whereOption) String() string {
	res := w.buildWhere(nil)
	bytes, _ := json.Marshal(res)
	return ldconv.BytesToStrUnsafe(bytes)
}

func (w *whereOption) buildGorm(db *GormDb) *GormDb {
	// log.Printf("=== db:%p", db)
	return w.where.buildGorm(db, w.value)
}

func (w *whereOption) buildWhere(db *GormDb) whereResult {
	return w.where.buildWhere(db, w.value)
}

func (w *whereOption) toTree() *whereOptionTree {
	return &whereOptionTree{
		Wheres: []whereOptionTreeNode{
			{Where: w},
		},
	}
}

func (w *whereOption) And(o interface{}) WhereOption {
	return w.toTree().And(o)
}

func (w *whereOption) Or(o interface{}) WhereOption {
	return w.toTree().Or(o)
}

type whereOptionTreeNode struct {
	Or    bool        `json:"or"`
	Where WhereOption `json:"where"`
}

type whereOptionTree struct {
	Wheres []whereOptionTreeNode `json:"wheres"`
}

func (w *whereOptionTree) String() string {
	res := w.buildWhere(nil)
	bytes, _ := json.Marshal(res)
	return ldconv.BytesToStrUnsafe(bytes)
}

func (w *whereOptionTree) buildWhere(db *GormDb) whereResult {
	res := w.Wheres[0].Where.buildWhere(db)
	if len(w.Wheres) == 1 {
		return res
	}

	res.Query = "(" + res.Query

	for _, v := range w.Wheres[1:] {
		tmp := v.Where.buildWhere(db)
		symbol := " AND "
		if v.Or {
			symbol = " OR "
		}

		res.Query = res.Query + symbol + tmp.Query
		res.Args = append(res.Args, tmp.Args...)
	}

	res.Query = res.Query + ")"
	return res
}

func (w *whereOptionTree) buildGorm(db *GormDb) *GormDb {
	// log.Printf("=== db:%p", db)
	res := w.buildWhere(db)
	if strings.HasPrefix(res.Query, "(") && strings.HasSuffix(res.Query, ")") {
		res.Query = res.Query[1 : len(res.Query)-1]
	}

	if res.IsValid() {
		db = db.Where(res.Query, res.Args...)
	}

	return db
}

func (w *whereOptionTree) clone() *whereOptionTree {
	c := *w
	return &c
}

func (w *whereOptionTree) And(where interface{}) WhereOption {
	return w.append(false, where)
}

func (w *whereOptionTree) Or(where interface{}) WhereOption {
	return w.append(true, where)
}

func (w *whereOptionTree) append(or bool, where interface{}) WhereOption {
	w = w.clone()
	w.Wheres = append(w.Wheres, whereOptionTreeNode{
		Or:    or,
		Where: Where(where),
	})

	return w
}

/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FieldOrderer interface {
	buildGorm(db *GormDb, field string) *GormDb

	Order(i int) FieldOrderer
	Desc() FieldOrderer
	Asc() FieldOrderer

	// ORDER BY FIELD(`{column}`,3,2,1,4)
	Field(fieldValues interface{}) FieldOrderer

	getOrder() int
}

type fieldOrder struct {
	OrderNum int           `json:"order"`
	IsDesc   bool          `json:"desc"`
	Fields   interface{}   `json:"fields"`
	fields   reflect.Value `json:"-"`
}

func FieldOrder(i int) FieldOrderer {
	return fieldOrder{
		OrderNum: i,
	}
}

func (that fieldOrder) getOrder() int { return that.OrderNum }

func (that fieldOrder) buildGorm(db *GormDb, field string) *GormDb {
	if that.fields.IsValid() {
		return that.buildGormWithField(db, field)
	}

	exp := fmt.Sprintf("`%s`", field)
	if that.IsDesc {
		exp = fmt.Sprintf("`%s` DESC", field)
	}

	db = db.Order(exp)
	return db
}

func (that fieldOrder) Order(i int) FieldOrderer {
	that.OrderNum = i
	return that
}

func (that fieldOrder) Desc() FieldOrderer {
	that.IsDesc = true
	return that
}

func (that fieldOrder) Asc() FieldOrderer {
	that.IsDesc = false
	return that
}

func (that fieldOrder) Field(fieldValues interface{}) FieldOrderer {
	if fieldValues == nil {
		that.Fields = nil
		that.fields = reflect.Value{}
		return that
	}

	v := reflect.ValueOf(fieldValues)
	switch v.Kind() {
	default:
		panic(fmt.Sprintf("the paramter of `ORDER BY FIELD` must be array or slice. %s",
			v.Type().String()))

	case reflect.Array:
		v = v.Slice(0, v.Len())
	case reflect.Slice:
		break
	}

	if v.Len() == 0 {
		panic(fmt.Sprintf("the paramter of `ORDER BY FIELD` must not be empty array"))
	}

	that.Fields = v.Interface()
	that.fields = v
	return that
}

func (that fieldOrder) buildGormWithField(db *GormDb, field string) *GormDb {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "FIELD(`%s`", field)

	l := that.fields.Len()
	args := make([]interface{}, 0, l)
	for i := 0; i < l; i++ {
		v := that.fields.Index(i)
		fmt.Fprintf(buf, ", ?")
		args = append(args, v.Interface())
	}

	fmt.Fprintf(buf, ")")
	if that.IsDesc {
		fmt.Fprintf(buf, " DESC")
	}
	// return db.Order(gorm.Expr(buf.String(), args...))
	return db.Clauses(clause.OrderBy{
		Expression: gorm.Expr(buf.String(), args...),
	})
}

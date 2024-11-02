/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldtagmap"
)

const _WHERE_TAG = "ldgormwhere"

var (
	_WHERE_FIELD_TYPE = reflect.TypeOf((*FieldWherer)(nil)).Elem()

	whereCache = &sync.Map{}
)

func BuildWhere(db *GormDb, where interface{}) *GormDb {
	if where == nil {
		return db
	}

	w := Where(where)
	return w.buildGorm(db)
}

type fieldWhereReflect struct {
	Tags       ldtagmap.Tags
	Name       string
	Order      int32
	FieldOrder int
	NotEmpty   bool
}

type whereReflect struct {
	Fields []*fieldWhereReflect
}

func (that *whereReflect) getTableName(val reflect.Value) string {
	if v, ok := val.Interface().(TableNamer); ok && v != nil {
		return v.TableName()
	}
	return ""
}

func (that *whereReflect) quote(db *GormDb, name string) string {
	if db == nil || db.Get() == nil || db.Get().Dialector == nil {
		if db != nil {
			if db.Get() != nil {
			}
		}
		return name
	}
	buf := &strings.Builder{}
	buf.Grow(len(name) * 2)
	db.Get().Dialector.QuoteTo(buf, name)
	return buf.String()
}

func (that *whereReflect) buildWhere(db *GormDb, val reflect.Value) whereResult {
	tableName := that.getTableName(val)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	buf := make([]string, 0, 2)
	wheres := make([]whereResult, 0, len(that.Fields))
	for _, f := range that.Fields {
		fw, ok := val.Field(f.FieldOrder).Interface().(FieldWherer)
		if !ok || fw.isEmpty() {
			if f.NotEmpty {
				panic(fmt.Sprintf("the where field must not be empty. %s", f.Name))
			}
			continue
		}

		buf = buf[:0]
		if tableName != "" {
			buf = append(buf, that.quote(db, tableName))
		}

		buf = append(buf, that.quote(db, f.Name))

		field := strings.Join(buf, ".")
		wheres = append(wheres, fw.buildWhere(field))
	}

	switch len(wheres) {
	case 0:
		return whereResult{}
	case 1:
		return wheres[0]
	}

	res := wheres[0]

	res.Query = "(" + res.Query
	for _, tmp := range wheres[1:] {
		res.Query = res.Query + " AND " + tmp.Query
		res.Args = append(res.Args, tmp.Args...)
	}
	res.Query = res.Query + ")"

	return res
}

func (that *whereReflect) buildGorm(db *GormDb, val reflect.Value) *GormDb {
	res := that.buildWhere(db, val)
	if strings.HasPrefix(res.Query, "(") && strings.HasSuffix(res.Query, ")") {
		res.Query = res.Query[1 : len(res.Query)-1]
	}

	if res.IsValid() {
		db = db.Where(res.Query, res.Args...)
	}

	return db
}

type sortSliceFieldInfo []*fieldWhereReflect

func (s sortSliceFieldInfo) Len() int      { return len(s) }
func (s sortSliceFieldInfo) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortSliceFieldInfo) Less(i, j int) bool {
	if s[i].Order != s[j].Order {
		return s[i].Order < s[j].Order
	}
	return s[i].FieldOrder < s[j].FieldOrder
}

func getWhereReflect(typ reflect.Type) *whereReflect {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("the where type must be `ldgorm.WhereOption` or struct or pointer to struct. %s", typ))
	}

	cache := whereCache
	if v, _ := cache.Load(typ); v != nil {
		tmp, ok := v.(*whereReflect)
		if ok {
			return tmp
		}
	}

	fields := make([]*fieldWhereReflect, 0, typ.NumField())
	for i, n := 0, typ.NumField(); i < n; i++ {
		f := getFieldWhereReflect(typ, i)
		if f == nil {
			continue
		}

		fields = append(fields, f)
	}

	if len(fields) == 0 {
		panic("where struct must have at least one where field")
	}

	sort.Sort(sortSliceFieldInfo(fields))
	w := &whereReflect{
		Fields: fields,
	}

	cache.Store(typ, w)
	return w
}

func getFieldWhereReflect(typ reflect.Type, i int) *fieldWhereReflect {
	field := typ.Field(i)
	tag, ok := field.Tag.Lookup(_WHERE_TAG)
	if !ok {
		return nil
	}
	if len(tag) == 0 {
		return nil
	}

	tags := ldtagmap.Parse(tag)
	if _, ok := tags["-"]; ok {
		return nil
	}

	name := tags.Get("column")
	if len(name) == 0 {
		name = tags.Get("name")
		if len(name) == 0 {
			return nil
		}
	}

	if !field.Type.Implements(_WHERE_FIELD_TYPE) {
		panic("where field type must be `ldgorm.FieldWherer`")
	}

	order := tags.Get("order")
	notEmpty := tags.Has("notempty")

	return &fieldWhereReflect{
		Tags:       tags,
		Name:       name,
		Order:      ldconv.AsInt32(order, math.MaxInt32),
		FieldOrder: i,
		NotEmpty:   notEmpty,
	}
}

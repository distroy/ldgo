/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/distroy/ldgo/ldconv"
	"github.com/distroy/ldgo/ldtagmap"
)

const _ORDER_TAG = "gormorder"

var (
	_ORDER_FIELD_TYPE = reflect.TypeOf((*FieldOrderer)(nil)).Elem()

	oderCache = &sync.Map{}
)

func Order(o interface{}) Option {
	if o == nil {
		return nil
	}

	val := reflect.ValueOf(o)
	ref := getOrderReflect(val.Type())

	return &orderOption{
		Value: val,
		Order: ref,
	}
}

type orderOption struct {
	Value reflect.Value
	Order *orderReflect
}

func (that *orderOption) String() string {
	bytes, _ := json.Marshal(that.Value.Interface())
	return ldconv.BytesToStrUnsafe(bytes)
}

func (that *orderOption) buildGorm(db *GormDb) *GormDb {
	return that.Order.buildOrder(db, that.Value)
}

type fieldOrderReflect struct {
	Tags       ldtagmap.Tags
	Name       string
	FieldOrder int
}

type orderReflect struct {
	Fields []*fieldOrderReflect
}

func (that *orderReflect) buildOrder(db *GormDb, val reflect.Value) *GormDb {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make([]fieldOrderTemp, 0, len(that.Fields))
	for _, fRef := range that.Fields {
		field, ok := val.Field(fRef.FieldOrder).Interface().(FieldOrderer)
		if !ok {
			continue
		}
		fields = append(fields, fieldOrderTemp{
			Order:  field,
			Relect: fRef,
		})
	}

	sort.Sort(sortSliceFieldOrderTemp(fields))

	for _, f := range fields {
		db = f.Order.buildGorm(db, f.Relect.Name)
	}

	return db
}

type fieldOrderTemp struct {
	Order  FieldOrderer
	Relect *fieldOrderReflect
}

type sortSliceFieldOrderTemp []fieldOrderTemp

func (s sortSliceFieldOrderTemp) Len() int      { return len(s) }
func (s sortSliceFieldOrderTemp) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortSliceFieldOrderTemp) Less(i, j int) bool {
	if s[i].Order.getOrder() != s[j].Order.getOrder() {
		return s[i].Order.getOrder() < s[j].Order.getOrder()
	}
	return s[i].Relect.FieldOrder < s[j].Relect.FieldOrder
}

func getOrderReflect(typ reflect.Type) *orderReflect {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("the order type must be struct or pointer to struct. %s", typ))
	}

	cache := oderCache
	if v, _ := cache.Load(typ); v != nil {
		tmp, ok := v.(*orderReflect)
		if ok {
			return tmp
		}
	}

	fields := make([]*fieldOrderReflect, 0, typ.NumField())
	for i, n := 0, typ.NumField(); i < n; i++ {
		f := getOrderFieldReflect(typ, i)
		if f == nil {
			continue
		}

		fields = append(fields, f)
	}

	if len(fields) == 0 {
		panic("struct must have at least one order field")
	}

	ref := &orderReflect{
		Fields: fields,
	}

	cache.Store(typ, ref)
	return ref
}

func getOrderFieldReflect(typ reflect.Type, i int) *fieldOrderReflect {
	field := typ.Field(i)
	tag, ok := field.Tag.Lookup(_ORDER_TAG)
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
		return nil
	}

	if !field.Type.Implements(_ORDER_FIELD_TYPE) {
		panic("order field type must be `ldgorm.FieldOrderer`")
	}

	return &fieldOrderReflect{
		Tags:       tags,
		Name:       name,
		FieldOrder: i,
	}
}

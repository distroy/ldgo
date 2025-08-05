/*
 * Copyright (C) distroy
 */

package jsontag

import (
	"reflect"
	"sync"
)

var (
	pool = &sync.Map{}
)

type Field struct {
	Field     reflect.StructField
	Type      reflect.Type // not be ptr
	Name      string
	Index     int
	OmitEmpty bool
}

type Struct struct {
	fields []Field
}

func (p *Struct) NumField() int      { return len(p.fields) }
func (p *Struct) Field(i int) *Field { return &p.fields[i] }

func Get(typ reflect.Type) *Struct {
	if i, ok := pool.Load(typ); ok {
		return i.(*Struct)
	}
	p := getStruct(typ)
	i, _ := pool.LoadOrStore(typ, p)
	return i.(*Struct)
}

func getStruct(typ reflect.Type) *Struct {
	n := typ.NumField()
	p := &Struct{
		fields: make([]Field, 0, n),
	}

	for i := range n {
		sf := typ.Field(i)
		f, ok := getField(i, sf)
		if ok {
			p.fields = append(p.fields, f)
		}
	}
	return p
}

func getField(idx int, sf reflect.StructField) (Field, bool) {
	ft := sf.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}

	// if !sf.IsExported() && !(sf.Anonymous && ft.Kind() == reflect.Struct) {
	// 	return Field{}, false
	// }
	if sf.Anonymous {
		if !sf.IsExported() && ft.Kind() != reflect.Struct {
			return Field{}, false
		}

	} else if !sf.IsExported() {
		return Field{}, false
	}

	tag := sf.Tag.Get("json")
	if tag == "-" {
		return Field{}, false
	}

	name, opts := parseTag(tag)
	if name == "" {
		name = sf.Name
	}

	return Field{
		Field:     sf,
		Type:      ft,
		Name:      name,
		Index:     idx,
		OmitEmpty: opts.Contains("omitempty"),
	}, true
}

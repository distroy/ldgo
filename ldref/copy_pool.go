/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strings"
	"sync"
)

const defaultCopyTagName = "json"

var copyTypePool = &sync.Map{}
var copyNameKeysInTag = []string{
	"name",
	"column",
}

type copyFieldInfo struct {
	reflect.StructField

	Name     string
	Index    int
	Ignore   bool
	TypeZero reflect.Value
}

type copyStructKey struct {
	Type    reflect.Type
	TagName string
}

type copyStructValue struct {
	Type    reflect.Type
	TagName string
	Fields  map[string]*copyFieldInfo
	Ignores []*copyFieldInfo
}

func getTagName(name string) string {
	if name != "" {
		return name
	}
	return defaultCopyTagName
}

func getCopyTypeInfo(typ reflect.Type, tagName string) *copyStructValue {
	tagName = getTagName(tagName)
	key := copyStructKey{
		Type:    typ,
		TagName: tagName,
	}
	i, _ := copyTypePool.Load(key)
	if i != nil {
		return i.(*copyStructValue)
	}

	p := parseCopyStructInfo(typ, tagName)
	copyTypePool.LoadOrStore(key, p)

	return p
}

func parseCopyStructInfo(typ reflect.Type, tagName string) *copyStructValue {
	num := typ.NumField()
	res := &copyStructValue{
		Type:   typ,
		Fields: make(map[string]*copyFieldInfo, num),
	}

	for i := 0; i < num; i++ {
		f := parseCopyFieldInfo(i, typ.Field(i), tagName)
		if f.Ignore {
			res.Ignores = append(res.Ignores, f)
			continue
		}
		res.Fields[f.Name] = f
	}

	return res
}

func parseCopyFieldInfo(index int, field reflect.StructField, tagName string) *copyFieldInfo {
	f := &copyFieldInfo{
		StructField: field,
		Name:        field.Name,
		Index:       index,
		TypeZero:    reflect.Zero(field.Type),
	}

	tagStr := field.Tag.Get(tagName)
	if tagStr == "" {
		return f
	}

	if tagStr == "-" {
		f.Ignore = true
		return f
	}

	tagList := strings.FieldsFunc(tagStr, func(r rune) bool { return r == ';' || r == ',' })

	name := tagList[0]
	if name == "" {
		return f
	}

	idx := strings.Index(name, ":")
	if idx < 0 {
		f.Name = name
		return f
	}

	tagMap := make(map[string]string, len(tagList))
	tagMap[strings.ToLower(name[:idx])] = name[idx+1:]
	for _, str := range tagList[1:] {
		idx := strings.Index(str, ":")
		if idx < 0 {
			// tagMap[str] = ""
			continue
		}
		tagMap[strings.ToLower(str[:idx])] = str[idx+1:]
	}

	for _, key := range copyNameKeysInTag {
		val := tagMap[key]
		if val != "" {
			f.Name = val
			return f
		}
	}

	return f
}

var copyFuncPool = &commFuncPool[copyFuncKey, copyFuncType]{}

type copyFuncKey struct {
	Target   copyStructKey
	Source   copyStructKey
	Clone    bool
	Indirect bool
}

func isBaseType(typ reflect.Type) bool {
	switch refKindOfType(typ) {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan:
		return isBaseType(typ.Elem())

	case reflect.Map:
		return isBaseType(typ.Elem()) && isBaseType(typ.Key())

	case reflect.Struct, reflect.Interface:
		return false
	}
	return true
}

func getCopyFunc(c *copyContext, tTyp, sTyp reflect.Type) (*copyFuncType, func()) {
	return _getCopyFuncWithIndirect(c, tTyp, sTyp, false)
}
func getCopyFuncIndirect(c *copyContext, tTyp, sTyp reflect.Type) (*copyFuncType, func()) {
	return _getCopyFuncWithIndirect(c, tTyp, sTyp, true)
}
func _getCopyFuncWithIndirect(c *copyContext, tTyp, sTyp reflect.Type, indirect bool) (*copyFuncType, func()) {
	key := copyFuncKey{
		Target:   copyStructKey{Type: tTyp},
		Source:   copyStructKey{Type: sTyp},
		Clone:    c.Clone,
		Indirect: indirect,
	}

	if !isBaseType(tTyp) {
		key.Target.TagName = getTagName(c.TargetTag)
	}
	if !isBaseType(sTyp) {
		key.Source.TagName = getTagName(c.SourceTag)
	}

	return copyFuncPool.Get(key, func() copyFuncType {
		if indirect {
			return _getCopyFuncReflectWithIndirect(c, tTyp, sTyp)
		}
		return _getCopyFuncReflect(c, tTyp, sTyp)
	})
}

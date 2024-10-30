/*
 * Copyright (C) distroy
 */

package ldsheet

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/distroy/ldgo/v3/ldtagmap"
)

type model struct {
	Inited  bool
	Headed  bool
	Type    reflect.Type
	Headers map[string]int
	Fields  []*modelField
}

type modelField struct {
	Field       reflect.StructField
	Tags        ldtagmap.Tags
	NotEmpty    bool
	IgnoreCase  bool
	Name        string
	Prefix      string
	FieldIndex  int
	HeaderIndex int
}

func parseModelField(field reflect.StructField) *modelField {
	tag, ok := field.Tag.Lookup(_TAG_NAME)
	if !ok {
		return &modelField{
			Field:       field,
			HeaderIndex: -1,
		}
	}

	tags := ldtagmap.Parse(tag)
	if _, ok := tags["-"]; ok {
		return nil
	}

	name := tags.Get("name")
	prefix := tags.Get("prefix")
	notEmpty := tags.Has("notempty")
	ignoreCase := tags.Has("ignorecase")

	return &modelField{
		Field:       field,
		Name:        name,
		Prefix:      prefix,
		NotEmpty:    notEmpty,
		IgnoreCase:  ignoreCase,
		HeaderIndex: -1,
	}
}

func (f *modelField) GetIndex(headers map[string]int) int {
	if f.IgnoreCase {
		temp := make(map[string]int, len(headers))
		for k, v := range headers {
			temp[strings.ToLower(k)] = v
		}
		headers = temp
	}

	if idx := f.getIndexByNameInTag(headers); idx >= 0 {
		return idx
	}

	if idx := f.getIndexByPrefix(headers); idx >= 0 {
		return idx
	}

	if f.Name != "" || f.Prefix != "" {
		return -1
	}

	if idx := f.getIndexByFieldName(headers); idx >= 0 {
		return idx
	}

	return -1
}

func (f *modelField) getIndexByNameInTag(headers map[string]int) int {
	name := f.Name
	name = f.getStr(name)
	if name != "" {
		idx, ok := headers[name]
		if ok {
			return idx
		}
	}

	return -1
}

func (f *modelField) getIndexByPrefix(headers map[string]int) int {
	prefix := f.Prefix
	prefix = f.getStr(prefix)
	if prefix != "" {
		for key, idx := range headers {
			key = f.getStr(key)
			if strings.HasPrefix(key, prefix) {
				return idx
			}
		}
	}

	return -1
}

func (f *modelField) getIndexByFieldName(headers map[string]int) int {
	name := f.Field.Name
	name = f.getStr(name)
	if idx, ok := headers[name]; ok {
		return idx
	}

	name = splitStringWord(f.Field.Name)
	if idx, ok := headers[name]; ok {
		return idx
	}

	name = strings.ToLower(name)
	if idx, ok := headers[name]; ok {
		return idx
	}

	return -1
}

func (f *modelField) getStr(s string) string {
	ignoreCase := f.IgnoreCase
	if !ignoreCase {
		return s
	}
	return strings.ToLower(s)
}

func (f *modelField) Validate() error {
	if f.HeaderIndex >= 0 {
		return nil
	}

	if !f.NotEmpty {
		return nil
	}

	field := f.Field.Name
	if len(f.Name) > 0 {
		return fmt.Errorf("[ldsheet] the field is missed in header. field:%s, name:%s", field, f.Name)
	}

	if len(f.Prefix) > 0 {
		return fmt.Errorf("[ldsheet] the field is missed in header. field:%s, prefix:%s", field, f.Prefix)
	}

	return fmt.Errorf("[ldsheet] the field is missed in header. field:%s", field)
}

func (f *modelField) ParseValue(obj reflect.Value, line []string) error {
	if f.HeaderIndex < 0 {
		return nil
	}

	if f.HeaderIndex >= len(line) || len(line[f.HeaderIndex]) == 0 {
		if f.NotEmpty {
			return fmt.Errorf("[ldsheet] the field must not be empty. field:%s", f.Field.Name)
		}
		return nil
	}

	str := line[f.HeaderIndex]
	fVal := obj.Field(f.FieldIndex)
	err := parseStringValue(fVal, str)
	if err != nil {
		return fmt.Errorf("[ldsheet] parse field value fail. type:%s, err:%s",
			fVal.Type().String(), err)
	}

	return nil
}

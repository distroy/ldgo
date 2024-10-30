/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"reflect"
	"strings"

	"github.com/distroy/ldgo/v3/ldconv"
)

func getReplaceSplits(splits []string) (l, r string) {
	switch len(splits) {
	case 0:
		return "{", "}"
	case 1:
		return splits[0], splits[0]
	default:
		return splits[0], splits[1]
	}
}

// StrMapReplace returns new string.
//
//	StrMapReplace("user: {user}", map[string]string{"user": "x"}) returns "user: x"
//	StrMapReplace("user: #user#", map[string]string{"user": "x"}, "#") returns "user: x"
//	StrMapReplace("user: <user>", map[string]string{"user": "x"}, "<", ">") returns "user: x"
func StrMapReplace(tmpl string, m map[string]string, splits ...string) string {
	if len(m) == 0 {
		return tmpl
	}

	l, r := getReplaceSplits(splits)

	replaceKeyValues := make([]string, 0, len(m)*2)
	for k, v := range m {
		kk := k
		vv := v

		kk = strings.Join([]string{l, kk, r}, "")
		if kk == "" {
			continue
		}
		replaceKeyValues = append(replaceKeyValues, kk, vv)
	}

	replacer := strings.NewReplacer(replaceKeyValues...)
	return replacer.Replace(tmpl)
}

// StrIMapReplace returns new string.
//
// m can be any map:
//
//	StrIMapReplace("user: {user}", map[interface{}]string{"user": "x"}) returns "user: x"
//	StrIMapReplace("user: #user#", map[interface{}]interface{}{"user": "x"}, "#") returns "user: x"
//	StrIMapReplace("count: <count>", map[string]int{"count": 123}, "<", ">") returns "count: 123"
func StrIMapReplace(tmpl string, m interface{}, splits ...string) string {
	if m, ok := m.(map[string]string); ok {
		return StrMapReplace(tmpl, m, splits...)
	}

	l, r := getReplaceSplits(splits)

	val := reflect.ValueOf(m)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Map {
		return ""
	}
	if val.Len() == 0 {
		return tmpl
	}

	replaceKeyValues := make([]string, 0, val.Len()*2)
	for it := val.MapRange(); it.Next(); {
		k := it.Key().Interface()
		v := it.Value().Interface()

		kk := ldconv.AsString(k)
		vv := ldconv.AsString(v)

		kk = strings.Join([]string{l, kk, r}, "")
		if kk == "" {
			continue
		}
		replaceKeyValues = append(replaceKeyValues, kk, vv)
	}

	replacer := strings.NewReplacer(replaceKeyValues...)
	return replacer.Replace(tmpl)
}

// StrMapParse use template to extract a string map form the text
func StrMapParse(tmpl, text string, splits ...string) (map[string]string, error) {
	l, r := getReplaceSplits(splits)
	parser := &StrMapParser{}

	err := parser.Init(tmpl, l, r)
	// log.Printf("%s", mustMarshalJson(parser.fields))
	if err != nil {
		return nil, err
	}

	res, err := parser.Parse(text)
	parser.Done()
	return res, err
}

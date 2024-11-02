/*
 * Copyright (C) distroy
 */

package ldsheet

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/distroy/ldgo/v3/ldconv"
)

var _CONV_FUNC_MAP = map[reflect.Kind]func(s string) (interface{}, error){
	reflect.String:  func(s string) (interface{}, error) { return s, nil },
	reflect.Bool:    func(s string) (interface{}, error) { return ldconv.ToBool(s) },
	reflect.Float32: func(s string) (interface{}, error) { return ldconv.ToFloat32(s) },
	reflect.Float64: func(s string) (interface{}, error) { return ldconv.ToFloat64(s) },
	reflect.Int:     func(s string) (interface{}, error) { return ldconv.ToInt(s) },
	reflect.Int8:    func(s string) (interface{}, error) { return ldconv.ToInt8(s) },
	reflect.Int16:   func(s string) (interface{}, error) { return ldconv.ToInt16(s) },
	reflect.Int32:   func(s string) (interface{}, error) { return ldconv.ToInt32(s) },
	reflect.Int64:   func(s string) (interface{}, error) { return ldconv.ToInt64(s) },
	reflect.Uint:    func(s string) (interface{}, error) { return ldconv.ToUint(s) },
	reflect.Uint8:   func(s string) (interface{}, error) { return ldconv.ToUint8(s) },
	reflect.Uint16:  func(s string) (interface{}, error) { return ldconv.ToUint16(s) },
	reflect.Uint32:  func(s string) (interface{}, error) { return ldconv.ToUint32(s) },
	reflect.Uint64:  func(s string) (interface{}, error) { return ldconv.ToUint64(s) },
}

func parseStringValue(val reflect.Value, str string) error {
	if val.Kind() == reflect.Ptr {
		val.Set(reflect.New(val.Type().Elem()))
		// return parseStringValue(val.Elem(), str)
		val = val.Elem()
	}

	conv := _CONV_FUNC_MAP[val.Kind()]
	if conv == nil {
		return fmt.Errorf("type has not been supported")
	}

	d, err := conv(str)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(d))
	return nil
}

func splitStringWord(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	res := make([]rune, 0, len(runes)*2)
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		if !unicode.IsUpper(curr) {
			res = append(res, curr)
			continue
		}

		if i > 0 {
			res = append(res, ' ')
		}

		last := curr
		j := i + 1
		for ; j < len(runes); j++ {
			curr := runes[j]
			if unicode.IsUpper(curr) {
				res = append(res, last)
				last = curr
				continue
			}

			if j > i+1 {
				res = append(res, ' ')
			}
			res = append(res, last)
			last = curr
			break
		}
		res = append(res, last)
		i = j
	}
	return string(res)
}

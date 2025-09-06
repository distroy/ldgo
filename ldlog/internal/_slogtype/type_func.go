/*
 * Copyright (C) distroy
 */

package _slogtype

import (
	"reflect"

	"github.com/distroy/ldgo/v3/ldlog/internal/_logref"
)

func toType[T, S any](v S) T                 { return _logref.ToType[T](v) }
func asType[T any](v any, def ...T) T        { return _logref.AsType(v, def...) }
func checkTypeEqual(this, that reflect.Type) { _logref.CheckTypeEqual(this, that) }

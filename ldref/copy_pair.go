/*
 * Copyright (C) distroy
 */

package ldref

import "reflect"

type copyPair struct {
	To   reflect.Kind
	From reflect.Kind
}

type (
	copyFuncType    = func(c *copyContext, target, source reflect.Value) (end bool)
	getCopyFuncType = func(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType
)

var (
	copyFuncMap    = map[copyPair]copyFuncType{}
	getCopyFuncMap = map[copyPair]getCopyFuncType{}
)

func registerCopyFunc(m map[copyPair]copyFuncType) {
	for pair, fnCopy := range m {
		copyFuncMap[pair] = fnCopy
	}
}

func registerGetCopyFunc(m map[copyPair]getCopyFuncType) {
	for pair, fn := range m {
		getCopyFuncMap[pair] = fn
	}
}

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		// {To: reflect.Ptr, From: reflect.Invalid}: copyReflectToPtrFromInvalid,
	})
}

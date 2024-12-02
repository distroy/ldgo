/*
 * Copyright (C) distroy
 */

package ldref

import "reflect"

var (
	typeOfByteSlice = reflect.TypeOf((*[]byte)(nil)).Elem()
	typeOfRuneSlice = reflect.TypeOf((*[]rune)(nil)).Elem()

	typeOfIface       = reflect.TypeOf((*interface{})(nil)).Elem()
	typeOfIfaceSlice  = reflect.TypeOf((*[]interface{})(nil)).Elem()
	typeOfMapStrIface = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
)

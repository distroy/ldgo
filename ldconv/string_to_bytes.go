/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"reflect"
	"unsafe"
)

func BytesToStr(b []byte) string {
	return string(b)
}

func StrToBytes(s string) []byte {
	return []byte(s)
}

// BytesToStrUnsafe the source bytes cannot free
func BytesToStrUnsafe(b []byte) string {
	// if b == nil {
	// 	return ""
	// }
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func bytesToStrUnsafeV1(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StrToBytesUnsafe the result bytes cannot be modified and the source string cannot free
func StrToBytesUnsafe(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
	// return *(*[]byte)(unsafe.Pointer(&s))
}

func strToBytesUnsafeV2(s string) []byte {
	type bytesFromString reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := &bytesFromString{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bh))
}

func strToBytesUnsafeV1(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

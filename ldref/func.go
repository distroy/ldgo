/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"runtime"
	"strings"
)

type FuncName struct {
	Full     string // github.com/distroy/ldgo/v3/ldref.(*testOject).Func
	Short    string // ldref.(*testOject).Func
	Path     string // github.com/distroy/ldgo/v3
	Package  string // ldref
	Method   string // Func
	Receiver string // *testOject
}

func GetFuncName(f interface{}) FuncName {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return ParseFuncName(fullName)
}

func ParseFuncName(fullName string) FuncName {
	shortName := fullName
	pathName := ""
	if idx := strings.LastIndex(fullName, "/"); idx >= 0 {
		pathName = fullName[:idx]
		shortName = fullName[idx+1:]
	}

	firstDotIndex := strings.Index(shortName, ".")

	packageName := shortName
	if firstDotIndex >= 0 {
		packageName = shortName[:firstDotIndex]
	}

	methodName := shortName
	receiverName := ""
	leftBracketIndex := strings.Index(shortName, "(")
	if leftBracketIndex < 0 {
		methodName = shortName[firstDotIndex+1:]

	} else {
		rightBracketIndex := strings.Index(shortName, ")")
		receiverName = shortName[leftBracketIndex+1 : rightBracketIndex]
		methodName = shortName[rightBracketIndex+2:]
	}

	return FuncName{
		Full:     fullName,
		Short:    shortName,
		Path:     pathName,
		Package:  packageName,
		Method:   methodName,
		Receiver: receiverName,
	}
}

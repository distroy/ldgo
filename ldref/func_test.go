/*
 * Copyright (C) distroy
 */

package ldref

import (
	"runtime"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type testFuncNameStruct struct{}

func (*testFuncNameStruct) Func() {}
func (*testFuncNameStruct) GetFuncName() string {
	f := func() string {
		f0 := func(f func() string) string {
			return f()
		}
		f1 := func() string {
			callerPc, _, _, _ := runtime.Caller(0)
			fullName := runtime.FuncForPC(callerPc).Name()
			return fullName
		}
		return f0(f1)
	}

	return f()
}

func TestGetFuncName(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("func", func(c convey.C) {
			f := GetFuncName(GetFuncName)
			c.So(f, convey.ShouldResemble, FuncName{
				Full:     "github.com/distroy/ldgo/v3/ldref.GetFuncName",
				Short:    "ldref.GetFuncName",
				Path:     "github.com/distroy/ldgo/v3",
				Package:  "ldref",
				Method:   "GetFuncName",
				Receiver: "",
			})
		})

		c.Convey("method", func(c convey.C) {
			f := GetFuncName((*testFuncNameStruct).Func)
			c.So(f, convey.ShouldResemble, FuncName{
				Full:     "github.com/distroy/ldgo/v3/ldref.(*testFuncNameStruct).Func",
				Short:    "ldref.(*testFuncNameStruct).Func",
				Path:     "github.com/distroy/ldgo/v3",
				Package:  "ldref",
				Method:   "Func",
				Receiver: "*testFuncNameStruct",
			})
		})

		c.Convey("method fm", func(c convey.C) {
			f := GetFuncName((&testFuncNameStruct{}).Func)
			c.So(f, convey.ShouldResemble, FuncName{
				Full:     "github.com/distroy/ldgo/v3/ldref.(*testFuncNameStruct).Func-fm",
				Short:    "ldref.(*testFuncNameStruct).Func-fm",
				Path:     "github.com/distroy/ldgo/v3",
				Package:  "ldref",
				Method:   "Func-fm",
				Receiver: "*testFuncNameStruct",
			})
		})

		c.Convey("testing.Main", func(c convey.C) {
			f := GetFuncName(testing.Main)
			c.So(f, convey.ShouldResemble, FuncName{
				Full:     "testing.Main",
				Short:    "testing.Main",
				Path:     "",
				Package:  "testing",
				Method:   "Main",
				Receiver: "",
			})
		})
	})
}

func TestParseFuncName(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("unamed func in method", func(c convey.C) {
			fullName := (&testFuncNameStruct{}).GetFuncName()
			r := ParseFuncName(fullName)
			c.So(r, convey.ShouldResemble, FuncName{
				Full:     "github.com/distroy/ldgo/v3/ldref.(*testFuncNameStruct).GetFuncName.func1.2",
				Short:    "ldref.(*testFuncNameStruct).GetFuncName.func1.2",
				Path:     "github.com/distroy/ldgo/v3",
				Package:  "ldref",
				Method:   "GetFuncName.func1.2",
				Receiver: "*testFuncNameStruct",
			})
		})

		c.Convey("unamed func in func", func(c convey.C) {
			callerPc, _, _, _ := runtime.Caller(0)
			fullName := runtime.FuncForPC(callerPc).Name()
			r := ParseFuncName(fullName)
			c.So(r, convey.ShouldResemble, FuncName{
				Full:     "github.com/distroy/ldgo/v3/ldref.TestParseFuncName.func1.2",
				Short:    "ldref.TestParseFuncName.func1.2",
				Path:     "github.com/distroy/ldgo/v3",
				Package:  "ldref",
				Method:   "TestParseFuncName.func1.2",
				Receiver: "",
			})
		})
	})
}

/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"fmt"
	"log"
	"reflect"
	"runtime"

	"github.com/distroy/ldgo/lderr"
	"github.com/gin-gonic/gin"
)

type (
	inConvType  = func(*Context) (reflect.Value, Error)
	outConvType = func(*Context, []reflect.Value)
)

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

type wrapperBase struct{}

func (w wrapperBase) isType(child, parent reflect.Type) bool {
	if child == parent {
		return true
	}
	if parent.Kind() == reflect.Interface && child.Implements(parent) {
		return true
	}
	return false
}

type wrapper struct {
	wrapperBase

	Method  string
	Path    string
	Handler string
	Name    string
	Type    reflect.Type
	InConvs []inConvType
	OutConv outConvType
}

func (w *wrapper) setMethod(method string) { w.Method = method }
func (w *wrapper) setPath(path string)     { w.Path = path }

func (w *wrapper) hasError(err Error) bool {
	return err != nil && err.Code() != 0
}

func (w *wrapper) returnError(c *Context, err Error) {
	c.AbortWithError(err)
}

func (w *wrapper) returnResponse(c *Context, rsp interface{}) {
	c.AbortWithData(rsp)
}

func (w *wrapper) getOutConv0() outConvType {
	return func(c *Context, v []reflect.Value) {}
}

func (w *wrapper) getOutConv1(outType reflect.Type) outConvType {
	errType := outType
	if !w.isType(errType, typeOfError) && !w.isType(errType, typeOfCommError) {
		panicf("%s output parameter type should be `ldgin.Error` or `error`", w.Name)
	}

	return func(c *Context, outs []reflect.Value) {
		out0 := outs[0].Interface()
		if err := out0; err != nil {
			if e := lderr.Wrap(err.(error)); w.hasError(e) {
				w.returnError(c, e)
				return
			}
		}
	}
}

func (w *wrapper) getAllInConvs(t reflect.Type) []inConvType {
	switch t.NumIn() {
	case 1:
		inType := t.In(0)
		if !w.isType(typeOfContext, inType) && !w.isType(typeOfGinContext, inType) {
			panicf("%s input parameter type should be `*ldgin.Context` or `*gin.Context`", w.Name)
		}

	case 2:
		inType := t.In(0)
		if !w.isType(typeOfContext, inType) && !w.isType(typeOfGinContext, inType) {
			panicf("%s first input parameter type should be `*ldgin.Context` or `*gin.Context`", w.Name)
		}

	default:
		panicf("%s input parameter count should be 1 or 2", w.Name)
	}

	ins := make([]inConvType, 0, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		t := t.In(i)
		conv := w.getInConv(t)
		ins = append(ins, conv)
	}

	return ins
}

func (w *wrapper) getInConv(t reflect.Type) inConvType {
	switch {
	case w.isType(typeOfContext, t):
		return func(c *Context) (reflect.Value, Error) {
			return reflect.ValueOf(c), nil
		}

	case w.isType(typeOfGinContext, t):
		return func(c *Context) (reflect.Value, Error) {
			return reflect.ValueOf(c.Gin()), nil
		}
	}

	convs := make([]func(*Context, reflect.Value) Error, 0, 2)
	convs = append(convs, w.getParserFunc(t))
	if f := w.getValidatorFunc(t); f != nil {
		convs = append(convs, f)
	}

	return func(c *Context) (reflect.Value, Error) {
		v := reflect.New(t.Elem())

		for _, f := range convs {
			if err := f(c, v); w.hasError(err) {
				return v, err
			}
		}

		c.Set(GinKeyRequest, v.Interface())
		return v, nil
	}
}

func (w *wrapper) getReqMethodByName(t reflect.Type, name string) func(*Context, reflect.Value) Error {
	m, ok := t.MethodByName(name)
	if !ok {
		return nil
	}

	mType := m.Type

	outNum := mType.NumOut()
	if outNum != 1 {
		log.Printf("output parameter count of request method should be 1. %s", m.Name)
		return nil
	}

	outType := mType.Out(0)
	if !w.isType(outType, typeOfError) && !w.isType(outType, typeOfCommError) {
		log.Printf("output parameter type of request method should be `ldgin.Error` or `error`. %s", m.Name)
		return nil
	}

	inNum := mType.NumIn()
	if inNum != 2 {
		log.Printf("input parameter count of request method should be 1. %s", m.Name)
		return nil
	}

	inType := mType.In(1)
	switch {
	default:
		log.Printf("input parameter type of request method should be `*ldgin.Context` or `*gin.Context`. %s", m.Name)
		return nil

	case w.isType(typeOfContext, inType):
		return func(c *Context, v reflect.Value) Error {
			ins := [2]reflect.Value{v, reflect.ValueOf(c)}
			outs := m.Func.Call(ins[:])
			err := outs[0].Interface()
			if err == nil {
				return nil
			}
			return lderr.Wrap(err.(error))
		}

	case w.isType(typeOfGinContext, inType):
		return func(c *Context, v reflect.Value) Error {
			ins := [2]reflect.Value{v, reflect.ValueOf(c.Gin())}
			outs := m.Func.Call(ins[:])
			err := outs[0].Interface()
			if err == nil {
				return nil
			}
			return lderr.Wrap(err.(error))
		}
	}
}

func (w *wrapper) getParserFunc(t reflect.Type) func(*Context, reflect.Value) Error {
	if f := w.getReqMethodByName(t, "Parse"); f != nil {
		return f
	}

	return func(c *Context, v reflect.Value) Error {
		return shouldBind(c, v.Interface())
	}
}

func (w *wrapper) getValidatorFunc(t reflect.Type) func(*Context, reflect.Value) Error {
	if f := w.getReqMethodByName(t, "Validate"); f != nil {
		return f
	}

	return nil
}

func (w *wrapper) call(g *gin.Context, h reflect.Value) {
	c := newCtxIfNotExists(g)
	if len(w.Handler) > 0 {
		c.setHandler(w.Handler)
	}
	if len(w.Method) > 0 {
		c.setMethod(w.Method)
	}
	if len(w.Path) > 0 {
		c.setPath(w.Path)
	}

	defer func() {
		if e := recover(); e != nil {
			seq := c.sequence

			array := [4 << 10]byte{}
			buf := array[:]
			buf = buf[:runtime.Stack(buf, false)]

			log.Printf("http: panic serving. remote=%s,sequence=%s,err=%v,stack=\n%s",
				c.Request.RemoteAddr, seq, e, buf)

			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
			w.returnError(c, lderr.Wrap(err, lderr.ErrServicePanic))
		}
	}()

	ins := make([]reflect.Value, 0, len(w.InConvs))
	for _, conv := range w.InConvs {
		v, err := conv(c)
		if w.hasError(err) {
			w.returnError(c, err)
			return
		}
		ins = append(ins, v)
	}

	outs := h.Call(ins)
	w.OutConv(c, outs)
}

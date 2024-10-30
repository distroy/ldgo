/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"log"
	"reflect"
	"runtime"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/gin-gonic/gin"
)

// WrapHandler returns then gin Handler
func WrapHandler(f Handler) func(c *gin.Context) {
	w := wrapHandler(f)
	return w.Do
}

// Wrap is the alias for WrapHandler
func Wrap(f Handler) func(c *gin.Context) { return WrapHandler(f) }

// WH is the alias for WrapHandler
func WH(f Handler) func(c *gin.Context) { return WrapHandler(f) }

type handler struct {
	wrapper

	Value reflect.Value
}

func wrapHandler(f Handler) *handler {
	handlerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	t := reflect.TypeOf(f)
	w := &handler{
		wrapper: wrapper{
			Handler: handlerName,
			Name:    "handler",
			Type:    t,
		},
		Value: reflect.ValueOf(f),
	}

	if t.Kind() != reflect.Func {
		panicf("%s should be a function", w.Name)
	}

	w.InConvs = w.getAllInConvs(t)

	switch t.NumOut() {
	case 0:
		w.OutConv = w.getOutConv0()

	case 1:
		outType := t.Out(0)
		w.OutConv = w.getOutConv1(outType)

	case 2:
		outTypes := []reflect.Type{t.Out(0), t.Out(1)}
		w.OutConv = w.getOutConv2(outTypes)

	default:
		panicf("%s output parameter count should be 0, 1 or 2", w.Name)
	}

	return w
}

func (w *handler) getOutConv1(outType reflect.Type) outConvType {
	errType := outType
	if !w.isType(errType, typeOfError) && !w.isType(errType, typeOfCommError) {
		panicf("%s output parameter type should be `ldgin.Error` or `error`", w.Name)
	}

	return func(c *Context, outs []reflect.Value) {
		out0 := outs[0].Interface()
		if err := out0; err != nil {
			if e := lderr.Wrap(err.(error)); w.hasError(e) {
				// w.returnError(c, e)
				c.AbortWithError(e)
				return
			}
		}

		// w.returnResponse(c, nil)
		c.AbortWithData(nil)
	}
}

func (w *handler) getOutConv2(outTypes []reflect.Type) outConvType {
	rspType := outTypes[0]
	errType := outTypes[1]

	if !w.isType(errType, typeOfError) && !w.isType(errType, typeOfCommError) {
		panicf("%s second output parameter type should be `ldgin.Error` or `error`", w.Name)
	}

	callRender := w.getRenderMethod(rspType)
	if callRender != nil {
		return func(c *Context, outs []reflect.Value) {
			out0 := outs[0].Interface()
			out1 := outs[1].Interface()

			if err := out1; err != nil {
				if e := lderr.Wrap(err.(error)); w.hasError(e) {
					// w.returnError(c, e)
					c.AbortWithError(e)
					return
				}
			}

			c.setRenderer(out0)
			callRender(c, outs[0])
		}
	}

	return func(c *Context, outs []reflect.Value) {
		out0 := outs[0].Interface()
		out1 := outs[1].Interface()

		if err := out1; err != nil {
			if e := lderr.Wrap(err.(error)); w.hasError(e) {
				// w.returnError(c, e)
				c.AbortWithErrorData(e, out0)
				return
			}
		}

		// w.returnResponse(c, out0)
		c.AbortWithData(out0)
	}
}

func (w handler) getRenderMethod(t reflect.Type) func(*Context, reflect.Value) {
	name := "Render"

	m, ok := t.MethodByName(name)
	if !ok {
		return nil
	}

	mType := m.Type

	outNum := mType.NumOut()
	if outNum != 0 {
		log.Printf("output parameter count of renderer method should be 0. %s", m.Name)
		return nil
	}

	inNum := mType.NumIn()
	if inNum != 2 {
		log.Printf("input parameter count of renderer method should be 1. %s", m.Name)
		return nil
	}

	inType := mType.In(1)
	switch {
	default:
		log.Printf("input parameter type of renderer method should be `*ldgin.Context` or `*gin.Context`. %s", m.Name)
		return nil

	case w.isType(typeOfContext, inType):
		return func(c *Context, v reflect.Value) {
			ins := [2]reflect.Value{v, reflect.ValueOf(c)}
			m.Func.Call(ins[:])
		}

	case w.isType(typeOfGinContext, inType):
		return func(c *Context, v reflect.Value) {
			ins := [2]reflect.Value{v, reflect.ValueOf(c.Gin())}
			m.Func.Call(ins[:])
		}
	}
}

func (w *handler) Do(g *gin.Context) {
	w.call(g, w.Value)
}

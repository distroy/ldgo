/*
 * Copyright (C) Shopee
 */

package ldgin

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

func Wrap(f Handler) func(c *gin.Context) {
	w := wrap(f)
	return w.Handler
}

type (
	inConvType  = func(Context, *gin.Context) (reflect.Value, Error)
	outConvType = func(Context, *gin.Context, []reflect.Value)
)

type wrapHandler struct {
	Type    reflect.Type
	Value   reflect.Value
	InConvs []inConvType
	outConv outConvType
}

func wrap(f Handler) *wrapHandler {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic("handler should be a function")
	}

	w := &wrapHandler{
		Type:  t,
		Value: reflect.ValueOf(f),
	}

	switch t.NumIn() {
	case 1:
		inType := t.In(0)
		if !inType.Implements(_TYPE_OF_CONTEXT) && !inType.Implements(_TYPE_OF_GIN_CONTEXT) {
			panic("handler input parameter type should be `ldgin.Context` or `*gin.Context`")
		}

	case 2:
		inType := t.In(0)
		if !inType.Implements(_TYPE_OF_CONTEXT) && !inType.Implements(_TYPE_OF_GIN_CONTEXT) {
			panic("handler first input parameter type should be `ldgin.Context` or `*gin.Context`")
		}

	default:
		panic("handler input parameter count should be 1 or 2")
	}

	ins := make([]inConvType, 0, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		t := t.In(i)
		conv := w.getInConv(t)
		ins = append(ins, conv)
	}
	w.InConvs = ins

	switch t.NumOut() {
	case 0:
		w.outConv = w.getOutConv0()
		break

	case 1:
		outType := t.Out(0)
		w.outConv = w.getOutConv1(outType)

	case 2:
		outTypes := []reflect.Type{t.Out(0), t.Out(1)}
		w.outConv = w.getOutConv2(outTypes)

	default:
		panic("handler output parameter count should be 0, 1 or 2")
	}

	return w
}

func (w *wrapHandler) getOutConv0() outConvType {
	return func(c1 Context, c2 *gin.Context, v []reflect.Value) {}
}

func (w *wrapHandler) returnError(c Context, g *gin.Context, err Error) {
	response := &commResponse{
		Sequence: GetSequence(c),
		Cost:     time.Since(GetBeginTime(c)).String(),
		ErrCode:  err.Code(),
		ErrMsg:   err.Error(),
	}
	g.Set(GIN_KEY_RESPONSE, response)
	g.JSON(err.Status(), response)
}

func (w *wrapHandler) returnResponse(c Context, g *gin.Context, rsp interface{}) {
	response := &commResponse{
		Sequence: GetSequence(c),
		Cost:     time.Since(GetBeginTime(c)).String(),
		Data:     rsp,
	}
	g.Set(GIN_KEY_RESPONSE, response)
	g.JSON(http.StatusOK, response)
}

func (w *wrapHandler) getOutConv1(outType reflect.Type) outConvType {
	if !outType.Implements(_TYPE_OF_ERROR) {
		panic("handler output parameter type should be `ldgin.Error`")
	}

	return func(c Context, g *gin.Context, outs []reflect.Value) {
		i := outs[0].Interface()
		err := i.(Error)
		if err != nil {
			w.returnError(c, g, err)
			return
		}

		w.returnResponse(c, g, nil)
	}
}

func (w *wrapHandler) getOutConv2(outTypes []reflect.Type) outConvType {
	rspType := outTypes[0]
	errType := outTypes[1]

	if !errType.Implements(_TYPE_OF_ERROR) {
		panic("handler second output parameter type should be `ldgin.Error`")
	}

	if rspType.Implements(_TYPE_OF_RENDERER) {
		return func(c Context, g *gin.Context, outs []reflect.Value) {
			err := outs[1].Interface().(Error)
			if err != nil {
				w.returnError(c, g, err)
				return
			}

			i := outs[0].Interface()
			render := i.(Renderer)
			render.Render(c)
		}
	}

	if rspType.Implements(_TYPE_OF_GIN_RENDERER) {
		return func(c Context, g *gin.Context, outs []reflect.Value) {
			err := outs[1].Interface().(Error)
			if err != nil {
				w.returnError(c, g, err)
				return
			}

			i := outs[0].Interface()
			render := i.(GinRenderer)
			render.Render(g)
		}
	}

	return func(c Context, g *gin.Context, outs []reflect.Value) {
		err := outs[1].Interface().(Error)
		if err != nil {
			w.returnError(c, g, err)
			return
		}

		rsp := outs[0].Interface()
		w.returnResponse(c, g, rsp)
	}
}

func (w *wrapHandler) getInConv(t reflect.Type) inConvType {
	switch {
	case t.Implements(_TYPE_OF_CONTEXT):
		return func(c Context, g *gin.Context) (reflect.Value, Error) {
			return reflect.ValueOf(c), nil
		}

	case t.Implements(_TYPE_OF_GIN_CONTEXT):
		return func(c Context, g *gin.Context) (reflect.Value, Error) {
			return reflect.ValueOf(g), nil
		}

	case t.Implements(_TYPE_OF_GIN_CONTEXT):
		return func(c Context, g *gin.Context) (reflect.Value, Error) {
			return reflect.ValueOf(g), nil
		}
	}

	convs := make([]func(Context, *gin.Context, reflect.Value) Error, 0, 2)
	convs = append(convs, w.getParserFunc(t))
	if f := w.getValidatorFunc(t); f != nil {
		convs = append(convs, f)
	}

	return func(c Context, g *gin.Context) (reflect.Value, Error) {
		v := reflect.New(t)
		for _, f := range convs {
			if err := f(c, g, v); err != nil {
				return v, err
			}
		}
		return v, nil
	}
}

func (w *wrapHandler) getParserFunc(t reflect.Type) func(Context, *gin.Context, reflect.Value) Error {
	if t.Implements(_TYPE_OF_PARSER) {
		return func(c Context, g *gin.Context, v reflect.Value) Error {
			paser, _ := v.Interface().(Parser)
			if err := paser.Parse(c); err != nil {
				return err
			}
			return nil
		}
	}

	if t.Implements(_TYPE_OF_GIN_PARSER) {
		return func(c Context, g *gin.Context, v reflect.Value) Error {
			paser, _ := v.Interface().(GinParser)
			if err := paser.Parse(g); err != nil {
				return err
			}
			return nil
		}
	}

	return func(c Context, g *gin.Context, v reflect.Value) Error {
		err := decodeHttpRequest(c, g, v.Interface())
		if err != nil {
			return commError{error: err}
		}
		return nil
	}
}

func (w *wrapHandler) getValidatorFunc(t reflect.Type) func(Context, *gin.Context, reflect.Value) Error {
	if t.Implements(_TYPE_OF_VALIDATER) {
		return func(c Context, g *gin.Context, v reflect.Value) Error {
			validator, _ := v.Interface().(Validator)
			if err := validator.Validate(c); err != nil {
				return err
			}

			return nil
		}
	}

	if t.Implements(_TYPE_OF_GIN_VALIDATER) {
		return func(c Context, g *gin.Context, v reflect.Value) Error {
			validator, _ := v.Interface().(GinValidator)
			if err := validator.Validate(g); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (w *wrapHandler) Handler(g *gin.Context) {
	c := GetContext(g)
	ins := make([]reflect.Value, 0, len(w.InConvs))
	for _, conv := range w.InConvs {
		v, err := conv(c, g)
		if err != nil {
			w.returnError(c, g, err)
			return
		}
		ins = append(ins, v)
	}

	outs := w.Value.Call(ins)
	if len(outs) == 0 {
		return
	}

	w.outConv(c, g, outs)
}

/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"reflect"
	"sync"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldref"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	requestTypes = &sync.Map{}
)

var requestBinds = []*requestBind{
	{Tag: "form", Func: wrapGinBindFunc((*gin.Context).ShouldBindQuery)},
	{Tag: "json", Func: wrapGinBindFunc((*gin.Context).ShouldBindJSON)},
	{Tag: "uri", Func: wrapGinBindFunc((*gin.Context).ShouldBindUri)},
	{Tag: "header", Func: wrapGinBindFunc((*gin.Context).ShouldBindHeader)},
	{Tag: "multipart", Func: parseMultipart},
}

type requestField struct {
	reflect.StructField

	Name string
}

type requestBind struct {
	Tag    string
	Func   func(c *Context, reqType *requestType, reqBind *requestBind, reqVal reflect.Value) error
	Fields []requestField
}

type requestType struct {
	Type    reflect.Type
	ReqPool sync.Pool
	ReqZero reflect.Value
	Binds   []*requestBind
}

func getRequestType(t reflect.Type) *requestType {
	if v, _ := requestTypes.Load(t); v != nil {
		reqT, ok := v.(*requestType)
		if ok {
			return reqT
		}
	}

	reqT := &requestType{
		Type:    t,
		ReqZero: reflect.Zero(t),
		Binds:   make([]*requestBind, 0, len(requestBinds)),
	}

	for _, v := range requestBinds {
		fields := getStructFieldsByTag(t, v.Tag)
		if len(fields) == 0 {
			continue
		}

		reqT.Binds = append(reqT.Binds, &requestBind{
			Tag:    v.Tag,
			Func:   v.Func,
			Fields: fields,
		})
	}

	if v, ok := requestTypes.LoadOrStore(t, reqT); ok {
		reqT, _ = v.(*requestType)
	}
	return reqT
}

func newRequest(reqType *requestType) reflect.Value {
	if i := reqType.ReqPool.Get(); i != nil {
		if v, ok := i.(reflect.Value); ok {
			return v
		}
	}

	return reflect.New(reqType.Type)
}

func delRequest(reqType *requestType, val reflect.Value) {
	val.Elem().Set(reqType.ReqZero)
	reqType.ReqPool.Put(val)
}

func getStructFieldsByTag(objT reflect.Type, tag string) []requestField {
	fields := make([]requestField, 0, objT.NumField())
	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		if field.Anonymous {
			continue
		}

		tagStr, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		if tagStr == "" || tagStr == "-" {
			continue
		}

		fields = append(fields, requestField{
			StructField: field,
			Name:        tagStr,
		})
	}

	return fields
}

func fillHttpRequestByFeilds(dst, src reflect.Value, fields []requestField) {
	for _, field := range fields {
		dst.Type().NumField()
		valSrc := src.FieldByIndex(field.Index)
		if ldref.IsValZero(valSrc) {
			continue
		}
		valDst := dst.FieldByIndex(field.Index)
		valDst.Set(valSrc)
	}
}

func wrapGinBindFunc(fn func(g *gin.Context, o interface{}) error) func(c *Context, reqType *requestType, reqBind *requestBind, reqVal reflect.Value) error {
	return func(c *Context, reqType *requestType, reqBind *requestBind, reqVal reflect.Value) error {
		reqNew := newRequest(reqType)
		if err := fn(c.Gin(), reqNew.Interface()); err != nil {
			ldctx.LogE(c, "ShouldBind() fail", zap.String("tag", reqBind.Tag), zap.Error(err))
			delRequest(reqType, reqNew)
			return lderr.ErrParseRequest
		}

		fillHttpRequestByFeilds(reqVal, reqNew.Elem(), reqBind.Fields)
		delRequest(reqType, reqNew)
		return nil
	}
}

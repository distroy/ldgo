/*
 * Copyright (C) Shopee
 */

package ldgin

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func decodeHttpRequest(ctx Context, c *gin.Context, req interface{}) error {
	reqV := reflect.ValueOf(req)
	if reqV.Kind() != reflect.Ptr {
		ctx.LogE("input request type must be pointer to struct", zap.String("type", reqV.Kind().String()))
		return fmt.Errorf("input request type must be pointer to struct")
	}
	reqV = reqV.Elem()
	if reqV.Kind() != reflect.Struct {
		ctx.LogE("input request type must be pointer to struct", zap.String("type", reqV.Kind().String()))
		return fmt.Errorf("input request type must be pointer to struct")
	}

	reqT := reqV.Type()

	if fields := getStructFieldsByTag(reqT, "form"); len(fields) != 0 {
		reqNew := reflect.New(reqV.Type())
		if err := c.ShouldBindQuery(reqNew.Interface()); err != nil {
			ctx.LogE("ShouldBindQuery() fail", zap.Error(err))
			return err
		}

		fillHttpRequestByFeilds(reqV, reqNew.Elem(), fields)
	}

	if fields := getStructFieldsByTag(reqT, "json"); len(fields) != 0 {
		reqNew := reflect.New(reqV.Type())
		if err := c.ShouldBindJSON(reqNew.Interface()); err != nil {
			ctx.LogE("ShouldBindJSON() fail", zap.Error(err))
			return err
		}

		fillHttpRequestByFeilds(reqV, reqNew.Elem(), fields)
	}

	if fields := getStructFieldsByTag(reqT, "uri"); len(fields) != 0 {
		reqNew := reflect.New(reqV.Type())
		if err := c.ShouldBindUri(reqNew.Interface()); err != nil {
			ctx.LogE("ShouldBindUri() fail", zap.Error(err))
			return err
		}

		fillHttpRequestByFeilds(reqV, reqNew.Elem(), fields)
	}

	return nil
}

func getStructFieldsByTag(objT reflect.Type, tag string) []reflect.StructField {
	fields := make([]reflect.StructField, 0, objT.NumField())
	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		tagStr, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		if tagStr == "" || tagStr == "-" {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}

func fillHttpRequestByFeilds(req, reqNew reflect.Value, fields []reflect.StructField) {
	for _, field := range fields {
		req.Type().NumField()
		val := req.FieldByIndex(field.Index)
		val.Set(reqNew.FieldByIndex(field.Index))
	}
}

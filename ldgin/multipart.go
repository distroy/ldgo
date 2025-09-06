/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"errors"
	"mime/multipart"
	"reflect"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldlog"
)

const (
	contentTypeMultipart = "multipart/form-data"
)

func parseMultipart(c *Context, reqType *requestType, reqBind *requestBind, reqVal reflect.Value) error {
	g := c.Gin()

	contentType := g.ContentType()
	if contentType != contentTypeMultipart {
		ldctx.LogW(c, "the content type should be `multipart/form-data`", ldlog.String("contentType", contentType))
		return nil
	}

	form, err := g.MultipartForm()
	if err != nil {
		ldctx.LogE(c, "ShouldBind() fail", ldlog.String("tag", reqBind.Tag), ldlog.Error(err))
		return lderr.ErrParseRequest
	}

	for _, field := range reqBind.Fields {
		value := reqVal.FieldByIndex(field.Index)

		name := field.Name
		files := form.File[name]

		if len(files) == 0 {
			continue
		}

		_, err := setByMultipartFormFile(value, field.StructField, files)
		if err != nil {
			ldctx.LogE(c, "ShouldBind() fail", ldlog.String("tag", reqBind.Tag), ldlog.Error(err))
			return lderr.ErrParseRequest
		}

		// if isSetted {
		// 	continue
		// }
	}
	return nil
}

func setByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSetted bool, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(), len(files), len(files))
		isSetted, err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSetted {
			return isSetted, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, errors.New("unsupported field type for multipart.FileHeader")
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSetted bool, err error) {
	if value.Len() != len(files) {
		return false, errors.New("unsupported len of array for []*multipart.FileHeader")
	}
	for i := range files {
		setted, err := setByMultipartFormFile(value.Index(i), field, files[i:i+1])
		if err != nil || !setted {
			return setted, err
		}
	}
	return true, nil
}

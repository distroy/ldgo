/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"reflect"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldlog"
)

// shouldBind will decode header/uri/json/query(form)
func shouldBind(c *Context, req interface{}) error {
	reqVal := reflect.ValueOf(req)
	if reqVal.Kind() != reflect.Ptr || reqVal.Elem().Kind() != reflect.Struct {
		ldctx.LogE(c, "input request type must be pointer to struct", ldlog.Stringer("type", reqVal.Kind()))
		return lderr.ErrInvalidRequestType
	}

	reqVal = reqVal.Elem()
	reqType := getRequestType(reqVal.Type())

	for _, reqBind := range reqType.Binds {
		if len(reqBind.Fields) == 0 {
			continue
		}

		err := reqBind.Func(c, reqType, reqBind, reqVal)
		if err != nil {
			return err
		}
	}

	return nil
}

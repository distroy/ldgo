/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"reflect"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/gin-gonic/gin"
)

type ctxKeyType int

const (
	ctxKeyContext ctxKeyType = iota
)

const (
	GinHeaderSequence = "x-ldgin-sequence"
	GinHeaderLatency  = "x-ldgin-latency"
)

const (
	GinKeyContext  = "x-ldgin-context"
	GinKeyError    = "x-ldgin-error"
	GinKeyRequest  = "x-ldgin-request"
	GinKeyRenderer = "x-ldgin-renderer"
	GinKeyResponse = "x-ldgin-response"
)

var (
	// typeOfInterface = reflect.TypeOf((*interface{})(nil)).Elem()

	typeOfGinHandlerFunc   = reflect.TypeOf((*gin.HandlerFunc)(nil)).Elem()
	typeOfGinHandlersChain = reflect.TypeOf((*gin.HandlersChain)(nil)).Elem()

	typeOfGinContext = reflect.TypeOf((*gin.Context)(nil))
	typeOfContext    = reflect.TypeOf((*Context)(nil))

	typeOfCommError = reflect.TypeOf((*error)(nil)).Elem()
	typeOfError     = reflect.TypeOf((*lderr.Error)(nil)).Elem()

	// typeOfParser         = reflect.TypeOf((*Parser)(nil)).Elem()
	// typeOfValidator      = reflect.TypeOf((*Validator)(nil)).Elem()
	// typeOfParseValidator = reflect.TypeOf((*ParseValidator)(nil)).Elem()

	// typeOfGinParser         = reflect.TypeOf((*GinParser)(nil)).Elem()
	// typeOfGinValidator      = reflect.TypeOf((*GinValidator)(nil)).Elem()
	// typeOfGinParseValidator = reflect.TypeOf((*GinParseValidator)(nil)).Elem()

	// typeOfRenderer    = reflect.TypeOf((*Renderer)(nil)).Elem()
	// typeOfGinRenderer = reflect.TypeOf((*GinRenderer)(nil)).Elem()
)

const (
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
)

const (
	crlf = "\r\n"
)

const (
	chunkedHeaderKey   = "Transfer-Encoding"
	chunkedHeaderValue = "chunked"
)

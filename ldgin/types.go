/*
 * Copyright (C) Shopee
 */

package ldgin

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	_TYPE_OF_GIN_CONTEXT = reflect.TypeOf((*gin.Context)(nil))
	_TYPE_OF_CONTEXT     = reflect.TypeOf((*Context)(nil)).Elem()

	_TYPE_OF_COMM_ERROR = reflect.TypeOf((*error)(nil)).Elem()
	_TYPE_OF_ERROR      = reflect.TypeOf((*Error)(nil)).Elem()

	_TYPE_OF_PARSER          = reflect.TypeOf((*Parser)(nil))
	_TYPE_OF_VALIDATER       = reflect.TypeOf((*Validator)(nil))
	_TYPE_OF_PARSE_VALIDATOR = reflect.TypeOf((*ParseValidator)(nil))

	_TYPE_OF_GIN_PARSER          = reflect.TypeOf((*GinParser)(nil))
	_TYPE_OF_GIN_VALIDATER       = reflect.TypeOf((*GinValidator)(nil))
	_TYPE_OF_GIN_PARSE_VALIDATOR = reflect.TypeOf((*GinParseValidator)(nil))

	_TYPE_OF_RENDERER = reflect.TypeOf((*Renderer)(nil))

	_TYPE_OF_GIN_RENDERER = reflect.TypeOf((*GinRenderer)(nil))
)

type Parser interface {
	Parse(Context) Error
}

type Validator interface {
	Validate(Context) Error
}

type ParseValidator interface {
	Parser
	Validator
}

type Renderer interface {
	Render(Context)
}

type GinParser interface {
	Parse(*gin.Context) Error
}

type GinValidator interface {
	Validate(*gin.Context) Error
}

type GinParseValidator interface {
	GinParser
	GinValidator
}

type GinRenderer interface {
	Render(*gin.Context)
}

// Request must be:
// ParseValidator
// Parser
// Validator
// GinParseValidator
// GinParser
// GinValidator
// interface{}
type Request interface{}

// Response must be:
// Renderer
// GinRenderer
// interface{}
type Response interface{}

// Handler must be:
// func (*gin.Context)
// func (*gin.Context, Request) Error
// func (*gin.Context) (Response, Error)
// func (*gin.Context, Request) (Response, Error)
// func (Context)
// func (Context, Request) Error
// func (Context) (Response, Error)
// func (Context, Request) (Response, Error)
type Handler interface{}

type Error interface {
	error
	Status() int
	Code() int
}

type commError struct {
	error
}

func (commError) Status() int { return http.StatusOK }
func (commError) Code() int   { return -1 }

type commResponse struct {
	ErrCode  int         `json:"code"`
	ErrMsg   string      `json:"msg"`
	Cost     string      `json:"cost"`
	Sequence string      `json:"sequence"`
	Data     interface{} `json:"data"`
}

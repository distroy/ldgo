/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"net/http"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/gin-gonic/gin"
)

type testRequest struct {
	ProjectId int64  `uri:"project_id"`
	ChannelId int64  `uri:"channel_id"`
	Page      int    `form:"page"`
	Where     string `json:"where"`
}

type testResponse struct {
	UserId int64 `json:"user_id"`
	ShopId int64 `json:"shop_id"`
}

type testGinParser testRequest

func (p *testGinParser) Parse(g *gin.Context) error {
	return shouldBind(GetContext(g), p)
}

type testGinParserFail struct{}

type testGinValidator testRequest

func (p *testGinValidator) Validate(g *gin.Context) lderr.Error {
	return nil
}

type testGinValidatorFail struct{}

func (p *testGinValidatorFail) Validate(g *gin.Context) lderr.Error {
	return lderr.ErrHttpReadBody
}

func (p *testGinParserFail) Parse(g *gin.Context) lderr.Error {
	return lderr.ErrHttpReadBody
}

type testGinRenderer struct {
	str string
}

func (p *testGinRenderer) Render(g *gin.Context) {
	g.String(http.StatusOK, p.str)
}

type testParser testRequest

func (p *testParser) Parse(c *Context) error {
	return shouldBind(c, p)
}

type testParserFail struct{}

func (p *testParserFail) Parse(c *Context) lderr.Error {
	return lderr.ErrHttpReadBody
}

type testValidator testRequest

func (p *testValidator) Validate(c *Context) lderr.Error {
	return nil
}

type testValidatorFail struct{}

func (p *testValidatorFail) Validate(c *Context) lderr.Error {
	return lderr.ErrHttpReadBody
}

type testParseValidator testRequest

func (p *testParseValidator) Parse(c *Context) error {
	return shouldBind(c, p)
}
func (p *testParseValidator) Validate(c *Context) lderr.Error {
	return nil
}

type testRenderer struct {
	str string
}

func (p *testRenderer) Render(c *Context) {
	c.Gin().String(http.StatusOK, p.str)
}

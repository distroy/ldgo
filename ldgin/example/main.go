/*
 * Copyright (C) distroy
 */

package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/distroy/ldgo/v2/ldconv"
	"github.com/distroy/ldgo/v2/ldctx"
	"github.com/distroy/ldgo/v2/lderr"
	"github.com/distroy/ldgo/v2/ldgin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	StdContext = context.Context
)

var (
	ErrTestOneError  = lderr.New(http.StatusOK, 1, "test one error 1")
	ErrTestDataError = lderr.New(http.StatusOK, 2, "test data error 1")
)

func testOneError(c *ldgin.Context) lderr.Error {
	ldctx.LogI(c, "", zap.String("method", c.GetMethod()), zap.String("path", c.GetPath()),
		zap.String("handler", c.GetHandler()))
	return ErrTestOneError
}

func testDataError(c *ldgin.Context) (string, lderr.Error) {
	return "abc", ErrTestDataError
}

type testBindReq struct {
	Uri      string `uri:"uri"`
	Query1   string `form:"query1"`
	Query2   int64  `form:"query2"`
	Language string `header:"accept-language"`
}

func testBind(ctx StdContext, req *testBindReq) (*testBindReq, lderr.Error) {
	g := ldgin.GetGin(ctx)
	c := ldgin.GetContext(g)
	ldctx.LogI(c, "", zap.String("method", c.GetMethod()), zap.String("path", c.GetPath()),
		zap.String("handler", c.GetHandler()))
	return req, nil
}

type testRenderer struct{}

func (_ *testRenderer) Render(c *gin.Context) {
	c.JSON(http.StatusOK, "abc")
}

func testRender(c *ldgin.Context) (*testRenderer, lderr.Error) {
	return &testRenderer{}, nil
}

type testValidateReq struct {
	Valid int64 `form:"valid"`
}

func (req *testValidateReq) Validate(c StdContext) lderr.Error {
	if req.Valid != 0 {
		return lderr.New(http.StatusOK, 111, fmt.Sprintf("invalid requet. valid=%v", req.Valid))
	}
	return nil
}

func testValidate(c *ldgin.Context, req *testValidateReq) lderr.Error {
	return nil
}

type testMultipartReq struct {
	Name string                `form:"name"`
	File *multipart.FileHeader `multipart:"file"`
}

type testMultipartRsp struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
}

func testMultipart(c *ldgin.Context, req *testMultipartReq) (*testMultipartRsp, lderr.Error) {
	name := req.Name

	file := req.File
	if file == nil {
		return nil, lderr.WithDetail(lderr.ErrParseRequest, "parse multipart fail")
	}

	if name == "" {
		name = file.Filename
	}

	f, err := file.Open()
	if err != nil {
		return nil, lderr.WithDetail(lderr.ErrParseRequest, "open multipart file fail")
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, lderr.WithDetail(lderr.ErrParseRequest, "read multipart file fail")
	}

	md5Sum := md5.Sum(buf)
	hash := hex.EncodeToString(md5Sum[:])
	return &testMultipartRsp{
		Hash: hash,
		Name: name,
	}, nil
}

type testParseReq struct {
	Query1 string `form:"query1"`
	Query2 int64  `form:"query2"`
}

func (req *testParseReq) Parse(c *ldgin.Context) lderr.Error {
	if err := c.ShouldBindQuery(req); err != nil {
		ldctx.LogE(c, "ShouldBindQuery() fail", zap.Error(err))
		return lderr.ErrParseRequest
	}
	return nil
}

func testParse(c *ldgin.Context, req *testParseReq) (*testParseReq, lderr.Error) {
	return req, nil
}

func testSucc(c *ldgin.Context) lderr.Error {
	return nil
}

func testPanic(c *ldgin.Context) lderr.Error {
	ldctx.LogI(c, "", zap.String("method", c.GetMethod()), zap.String("path", c.GetPath()),
		zap.String("handler", c.GetHandler()))

	var p *int
	*p = 1
	return nil
}

type testReaderForChunked struct {
	Error error
	Lines []string

	index int
}

func (r *testReaderForChunked) Read(buff []byte) (int, error) {
	if r.index < len(r.Lines) {
		data := r.Lines[r.index]
		r.index++

		return copy(buff, data), nil
	}

	err := r.Error
	if err != nil {
		return 0, err
	}
	return 0, io.EOF
}

func testChunked(c *ldgin.Context) (ldgin.ReaderRenderer, lderr.Error) {
	if ldconv.AsBool(c.Query("download")) {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "chunked.text"))
	}

	var err error
	if ldconv.AsBool(c.Query("error")) {
		err = fmt.Errorf("test error")
	}

	// c.Header("Transfer-Encoding", "chunked")
	return ldgin.ReaderRenderer{
		ContentType: "text/plain; charset=utf-8",
		Reader: &testReaderForChunked{
			Error: err,
			Lines: []string{
				strings.Repeat("abc\n", 3),
				strings.Repeat("xyz\n", 5),
			},
		},
	}, nil
}

func initRouter(c ldctx.Context, router gin.IRouter) {
	r := ldgin.WrapGin(router)

	// r = r.WithAppPath("aaaaaa/")
	r = r.Group("/test", midware1)
	r.GET("panic", testPanic)

	r = r.Use(midware2)
	r.GET("/one_error", testOneError)
	r.GET("/data_error", testDataError)
	r.GET("/bind/:uri", testBind)
	r.GET("/render/", testRender)
	r.GET("/validate", testValidate)
	r.GET("/parse", testParse)
	r.GET("/midware_error", testSucc, midware3)
	r.POST("/multipart", testMultipart)
	r.GET("/chunked", testChunked)

	r.Static(`static`, ".")
	r.StaticFS(`static_fs`, ldgin.Dir(".", true))
	r.StaticFile(`static_file`, "main.go")
}

func main() {
	ctx := ldctx.Default()

	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	initRouter(ctx, engine)

	server := &http.Server{
		Addr:    ":8000",
		Handler: engine,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe() fail. err:%v", err)
	}
}

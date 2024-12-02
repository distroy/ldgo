/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/gin-gonic/gin"
)

func emptyMidwareFunc(g *gin.Context) {}

type Reporter interface {
	Report(g *gin.Context, cost time.Duration, url string, err error)
}

func ReportMidware(r Reporter) func(g *gin.Context) {
	if r == nil {
		return emptyMidwareFunc
	}
	m := &reportMidware{reporter: r}
	return m.do
}

type reportMidware struct {
	reporter Reporter
}

func (m reportMidware) do(g *gin.Context) {
	c := GetContext(g)

	start := c.GetBeginTime()

	httpReqMethod := c.Request.Method
	httpReqPath := c.Request.URL.Path

	c.Next()

	// 获得接口路径和错误码
	httpCode := c.Writer.Status()
	method := c.GetMethod()
	if method == "" {
		method = httpReqMethod
	}
	path := c.GetPath()
	if path == "" {
		path = httpReqPath
	}

	url := fmt.Sprintf("[%s] %s", method, path)

	// 计算耗时
	cost := time.Since(start)

	// 获取业务的错误码
	bizCode := 0
	errMsg := ""
	if rsp := GetResponse(c); rsp != nil {
		bizCode = rsp.Error.Code
		errMsg = rsp.Error.Message
	}

	if err := c.GetError(); lderr.IsSuccess(err) {
		bizCode = lderr.GetCode(err)
		errMsg = lderr.GetMessage(err)
	}

	if bizCode == 0 && httpCode != http.StatusOK {
		err := lderr.ErrUnkown
		bizCode = err.Code()
		errMsg = err.Error()
	}

	// 监控上报
	if r := m.reporter; r != nil {
		err := lderr.New(httpCode, bizCode, errMsg)
		r.Report(g, cost, url, err)
	}
}

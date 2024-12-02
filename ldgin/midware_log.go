/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"net/http"
	"time"

	"github.com/distroy/ldgo/v2/ldctx"
	"github.com/distroy/ldgo/v2/lderr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogMidware() func(g *gin.Context) { return logMidwareFunc }

func logMidwareFunc(g *gin.Context) {
	c := GetContext(g)

	start := c.GetBeginTime()

	httpReqMethod := c.Request.Method
	httpReqPath := c.Request.URL.Path

	l := ldctx.GetLogger(g)
	l = l.With(zap.String("method", httpReqMethod), zap.String("path", httpReqPath))
	l.Info("http request begin")

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

	// 计算耗时
	cost := time.Since(start)

	reqField := zap.Skip()
	if req := GetRequest(c); req != nil {
		reqField = zap.Reflect("req", req)
	}

	// 获取业务的错误码
	bizCode := 0
	errMsg := ""
	rspDataField := zap.Skip()
	if rsp := GetResponse(c); rsp != nil {
		bizCode = rsp.Error.Code
		errMsg = rsp.Error.Message
		rspDataField = zap.Reflect("rspData", rsp.Data)
	}

	if err := c.GetError(); !lderr.IsSuccess(err) {
		bizCode = lderr.GetCode(err)
		errMsg = lderr.GetMessage(err)
	}

	if bizCode == 0 && httpCode != http.StatusOK {
		err := lderr.ErrUnkown
		bizCode = err.Code()
		errMsg = err.Error()
	}

	l.Info("http request end", zap.Int("httpCode", httpCode), zap.Int("bizCode", bizCode),
		zap.String("errmsg", errMsg), zap.Duration("cost", cost), reqField, rspDataField)
}

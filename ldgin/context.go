/*
 * Copyright (C) Shopee
 */

package ldgin

import (
	"encoding/hex"
	"time"

	"github.com/distroy/ldgo/ldcontext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKeyType int

const (
	ctxKey_GinContext ctxKeyType = iota
	ctxKey_Sequence
	ctxKey_BeginTime
)

const (
	GIN_KEY_CONTEXT  = "x-gin-context"
	GIN_KEY_RESPONSE = "x-gin-response"
)

type (
	ginContext = *gin.Context
	ldContext  = ldcontext.Context
)

type Context interface {
	ldcontext.Context
}

func newContext(g *gin.Context) Context {
	c := ldcontext.NewContext(g)

	c = c.WithValue(ctxKey_BeginTime, time.Now())
	c = c.WithValue(ctxKey_GinContext, g)

	uuid, _ := uuid.NewRandom()
	seq := hex.EncodeToString(uuid[:])
	c = c.With(zap.String("sequence", seq))
	c = c.WithValue(ctxKey_Sequence, seq)
	g.Header("x-gin-sequence", seq)

	g.Set(GIN_KEY_CONTEXT, c)

	return c
}

func GetContext(g *gin.Context) Context {
	v := g.Value(GIN_KEY_CONTEXT)
	c, _ := v.(ldContext)
	if c == nil {
		c = newContext(g)
	}
	return c
}

func GetGin(c Context) *gin.Context {
	v := c.Value(ctxKey_GinContext)
	t, ok := v.(*gin.Context)
	if !ok {
		return nil
	}
	return t
}

func GetResposeByGin(g *gin.Context) *commResponse {
	v := g.Value(GIN_KEY_CONTEXT)
	rsp, _ := v.(*commResponse)
	return rsp
}

func GetRespose(c Context) *commResponse {
	return GetResposeByGin(GetGin(c))
}

func GetBeginTimeByGin(g *gin.Context) time.Time {
	return GetBeginTime(GetContext(g))
}

func GetBeginTime(c Context) time.Time {
	v := c.Value(ctxKey_BeginTime)
	t, ok := v.(time.Time)
	if !ok {
		return time.Time{}
	}
	return t
}

func GetSequenceByGin(g *gin.Context) string {
	return GetSequence(GetContext(g))
}

func GetSequence(c Context) string {

	v := c.Value(ctxKey_Sequence)
	t, ok := v.(string)
	if !ok {
		return ""
	}
	return t
}

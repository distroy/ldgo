/*
 * Copyright (C) distroy
 */

package internal

import (
	"context"

	"github.com/distroy/ldgo/v3/ldctx"
)

type ctxKey int

const (
	ctxKeyInProcess ctxKey = iota
)

func InProcess(c context.Context) bool {
	b, _ := c.Value(ctxKeyInProcess).(bool)
	return b
}

func NewContext(c context.Context) context.Context {
	return ldctx.WithValue(c, ctxKeyInProcess, true)
}

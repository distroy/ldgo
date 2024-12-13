/*
 * Copyright (C) distroy
 */

package main

import (
	"net/http"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldgin"
)

func midware1(c *ldgin.Context) {
	ldctx.LogI(c, "midware1")
}

func midware2(c *ldgin.Context) lderr.Error {
	ldctx.LogI(c, "midware2")
	return nil
}

func midware3(c *ldgin.Context) lderr.Error {
	ldctx.LogI(c, "midware3")
	return lderr.New(http.StatusOK, 120, "midware error")
}

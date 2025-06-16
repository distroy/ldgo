/*
 * Copyright (C) distroy
 */

package handler

import "github.com/distroy/ldgo/v3/ldlog/internal/buffer"

type Buffer = buffer.Buffer

func newBuffer() *buffer.Buffer { return buffer.NewBuffer() }

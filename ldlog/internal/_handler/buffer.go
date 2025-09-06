/*
 * Copyright (C) distroy
 */

package _handler

import (
	"github.com/distroy/ldgo/v3/ldlog/internal/_buffer"
)

type Buffer = _buffer.Buffer

func newBuffer() *Buffer { return _buffer.NewBuffer() }

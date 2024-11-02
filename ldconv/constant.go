/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"math/big"

	"github.com/distroy/ldgo/v3/lderr"
)

var (
	_ERR_INVALID_SYNTAX = lderr.ErrInvalidNumberSyntax
	_ERR_INVALID_TYPE   = lderr.ErrInvalidConvertType

	_BIG_FLOAT_ZERO = (&big.Float{}).SetInt64(0)
)

func newBigFloatZero() *big.Float {
	return (&big.Float{}).Copy(_BIG_FLOAT_ZERO)
}

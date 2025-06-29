/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"os"

	"github.com/distroy/ldgo/v3/ldatomic"
)

var (
	defLogger = ldatomic.NewAny(New(os.Stderr))
	console   = New(os.Stderr)
	discard   = newDiscard()
)

func SetDefault(l *Logger) { defLogger.Store(l) }

func Default() *Logger { return defLogger.Load() }
func Console() *Logger { return console }
func Discard() *Logger { return discard }

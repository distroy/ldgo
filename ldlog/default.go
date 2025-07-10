/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"os"

	"github.com/distroy/ldgo/v3/ldatomic"
)

var (
	defLogger = ldatomic.NewAny(New(NewHandler(os.Stderr, nil)))
	console   = New(NewHandler(os.Stderr, nil))
	discard   = newDiscard()
)

func SetDefault(l *Logger) { defLogger.Store(l) }

func Default() *Logger { return defLogger.Load() }
func Console() *Logger { return console }
func Discard() *Logger { return discard }

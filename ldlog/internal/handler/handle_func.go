/*
 * Copyright (C) distroy
 */

package handler

import (
	"log/slog"
	_ "unsafe"
)

//go:linkname needsQuoting log/slog.needsQuoting
func needsQuoting(s string) bool

//go:linkname appendTextValue log/slog.appendTextValue
func appendTextValue(s *handleState, v slog.Value) error

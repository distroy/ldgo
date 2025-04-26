/*
 * Copyright (C) distroy
 */

package handler

import (
	"log/slog"
	"time"
	_ "unsafe"
)

//go:linkname appendRFC3339Millis log/slog.appendRFC3339Millis
func appendRFC3339Millis(b []byte, t time.Time) []byte

//go:linkname appendJSONTime log/slog.appendJSONTime
func appendJSONTime(s *handleState, t time.Time)

//go:linkname appendEscapedJSONString log/slog.appendEscapedJSONString
func appendEscapedJSONString(buf []byte, s string) []byte

//go:linkname needsQuoting log/slog.needsQuoting
func needsQuoting(s string) bool

//go:linkname appendJSONValue log/slog.appendJSONValue
func appendJSONValue(s *handleState, v slog.Value) error

//go:linkname appendTextValue log/slog.appendTextValue
func appendTextValue(s *handleState, v slog.Value) error

/*
 * Copyright (C) distroy
 */

package handler

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

func init() {
	checkTypeEqual(reflect.TypeOf(Source{}), reflect.TypeOf(slog.Source{}))
}

func GetSourcePtr(v *slog.Source) *Source { return toType[*Source](v) }
func GetSource(v slog.Source) Source      { return *GetSourcePtr(&v) }

// Source describes the location of a line of source code.
type Source struct {
	// Function is the package path-qualified function name containing the
	// source line. If non-empty, this string uniquely identifies a single
	// function in the program. This may be the empty string if not known.
	Function string `json:"function"`
	// File and Line are the file name and line number (1-based) of the source
	// line. These may be the empty string and zero, respectively, if not known.
	File string `json:"file"`
	Line int    `json:"line"`
}

func (s *Source) Get() *slog.Source { return toType[*slog.Source](s) }

func (s *Source) Group() Value {
	var as []Attr
	caller := s.Caller()
	if caller != "" {
		as = append(as, String("caller", caller))
	}
	// if s.Function != "" {
	// 	as = append(as, String("function", s.Function))
	// }
	// if s.File != "" {
	// 	as = append(as, String("file", s.File))
	// }
	// if s.Line != 0 {
	// 	as = append(as, Int("line", s.Line))
	// }
	return GroupValue(as...)
}

func (s *Source) Caller() string {
	file := s.File
	line := s.Line
	if file == "" {
		return ""
	}
	file = s.shortFilePath(file)
	caller := fmt.Sprintf("%s:%d", file, line)
	return caller
}

func (s *Source) shortFilePath(file string) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return file
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx < 0 {
		return file
	}
	return file[idx+1:]
}

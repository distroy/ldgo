/*
 * Copyright (C) distroy
 */

package handler

import (
	"log/slog"
	"reflect"
	"runtime"
	"time"
)

func init() {
	checkTypeEqual(reflect.TypeOf(Record{}), reflect.TypeOf(slog.Record{}))
}

func GetRecordPtr(v *slog.Record) *Record { return toType[*Record]((v)) }
func GetRecord(v slog.Record) Record      { return *GetRecordPtr(&v) }

const nAttrsInline = 5

type Level = slog.Level

type Record struct {
	// The time at which the output method (Log, Info, etc.) was called.
	Time time.Time

	// The log message.
	Message string

	// The level of the event.
	Level Level

	// The program counter at the time the record was constructed, as determined
	// by runtime.Callers. If zero, no program counter is available.
	//
	// The only valid use for this value is as an argument to
	// [runtime.CallersFrames]. In particular, it must not be passed to
	// [runtime.FuncForPC].
	PC uintptr

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Attrs.
	front [nAttrsInline]Attr

	// The number of Attrs in front.
	nFront int

	// The list of Attrs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []Attr
}

func (r *Record) Get() *slog.Record { return toType[*slog.Record](r) }

func (r *Record) Clone() Record { return GetRecord(r.Get().Clone()) }

func (r *Record) NumAttrs() int { return r.Get().NumAttrs() }
func (r *Record) Attrs(f func(Attr) bool) {
	r.Get().Attrs(func(a slog.Attr) bool { return f(GetAttr(a)) })
}

func (r *Record) AddAttrs(attrs ...Attr) { r.Get().AddAttrs(GetSAttrs(attrs)...) }
func (r *Record) Add(args ...any)        { r.Get().Add(args...) }

func (r *Record) Source() *Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return &Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

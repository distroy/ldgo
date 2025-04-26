/*
 * Copyright (C) distroy
 */

package handler

import (
	"log/slog"
	"unsafe"
)

func getAttr(v *slog.Attr) *Attr       { return (*Attr)(unsafe.Pointer(v)) }
func getValue(v *slog.Value) *Value    { return (*Value)(unsafe.Pointer(v)) }
func getSource(v *slog.Source) *Source { return (*Source)(unsafe.Pointer(v)) }
func getRecord(v *slog.Record) *Record { return (*Record)(unsafe.Pointer(v)) }

type (
	Attr   slog.Attr
	Source slog.Source
	Value  slog.Value
	Record slog.Record
)

//go:linkname isAttrEmpty log/slog.(Attr).isEmpty
func isAttrEmpty(a Attr) bool
func (a Attr) isEmpty() bool { return isAttrEmpty(a) }

//go:linkname countEmptyGroups log/slog.countEmptyGroups
func countEmptyGroups(as []slog.Attr) int

//go:linkname getSourceGroup log/slog.(*Source).group
func getSourceGroup(s *Source) slog.Value
func (s *Source) group() slog.Value { return getSourceGroup(s) }

//go:linkname isValueEmptyGroup log/slog.(*Value).isEmptyGroup
func isValueEmptyGroup(v *Value) bool
func (v *Value) isEmptyGroup() bool { return isValueEmptyGroup(v) }

//go:linkname getRecordSource log/slog.(*Record).source
func getRecordSource(r Record) *slog.Source
func (r Record) source() *slog.Source { return getRecordSource(r) }

type value struct {
	_ [0]func() // disallow ==
	// num holds the value for Kinds Int64, Uint64, Float64, Bool and Duration,
	// the string length for KindString, and nanoseconds since the epoch for KindTime.
	num uint64
	// If any is of type Kind, then the value is in num as described above.
	// If any is of type *time.Location, then the Kind is Time and time.Time value
	// can be constructed from the Unix nanos in num and the location (monotonic time
	// is not preserved).
	// If any is of type stringptr, then the Kind is String and the string value
	// consists of the length in num and the pointer in any.
	// Otherwise, the Kind is Any and any is the value.
	// (This implies that Attrs cannot store values of type Kind, *time.Location
	// or stringptr.)
	any any
}

func toValue(v *slog.Value) *value { return (*value)(unsafe.Pointer(v)) }

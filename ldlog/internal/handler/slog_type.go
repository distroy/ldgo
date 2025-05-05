/*
 * Copyright (C) distroy
 */

package handler

import (
	"fmt"
	"log/slog"
	"reflect"
	"unsafe"
)

func toType[T, S any](v S) T { return *(*T)(unsafe.Pointer(&v)) }

func getAttr(v *slog.Attr) *Attr       { return toType[*Attr](v) }
func getValue(v *slog.Value) *Value    { return toType[*Value](v) }
func getSource(v *slog.Source) *Source { return toType[*Source](v) }
func getRecord(v *slog.Record) *Record { return toType[*Record]((v)) }

type (
	Attr   slog.Attr
	Source slog.Source
	Record slog.Record
)

//go:linkname isAttrEmpty log/slog.(*Attr).isEmpty
func isAttrEmpty(a Attr) bool
func (a *Attr) isEmpty() bool { return isAttrEmpty(*a) }

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
func (r *Record) source() *slog.Source { return getRecordSource(*r) }

type Value struct {
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

func init() {
	typePairs := [][2]any{
		{Value{}, slog.Value{}},
		{Attr{}, slog.Attr{}},
		{Source{}, slog.Source{}},
		{Record{}, slog.Record{}},
	}
	for _, pair := range typePairs {
		v0, v1 := pair[0], pair[1]
		if t0, t1 := reflect.TypeOf(v0), reflect.TypeOf(v1); !isTypeEqual(t0, t1) {
			panic(fmt.Errorf("%s not not compatible with %s", t0.String(), t1.String()))
		}
	}
}

func isTypeEqual(t0, t1 reflect.Type) bool {
	if t0 == t1 {
		return true
	}
	if t0.Kind() != t1.Kind() {
		return false
	}

	switch t0.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.UnsafePointer:
		return true

	case reflect.Array:
		return t0.Len() == t1.Len() && isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Slice, reflect.Chan, reflect.Pointer:
		return isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Map:
		return isTypeEqual(t0.Key(), t1.Key()) && isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Func:
		return isTypeEqualForFunc(t0, t1)

	case reflect.Struct:
		return isTypeEqualForStruct(t0, t1)

	case reflect.Interface:
		return t0.Implements(t1) && t1.Implements(t0)
	}
	return false
}

func isTypeEqualForFunc(t0, t1 reflect.Type) bool {
	nout0 := t0.NumOut()
	nout1 := t1.NumOut()
	nin0 := t0.NumIn()
	nin1 := t1.NumIn()
	if nout0 != nout1 || nin0 != nin1 {
		return false
	}
	for i := range nout0 {
		in0 := t0.In(i)
		in1 := t1.In(i)
		if !isTypeEqual(in0, in1) {
			return false
		}
	}
	for i := range nin0 {
		out0 := t0.Out(i)
		out1 := t1.Out(i)
		if !isTypeEqual(out0, out1) {
			return false
		}
	}
	return true
}

func isTypeEqualForStruct(t0, t1 reflect.Type) bool {
	n0 := t0.NumField()
	n1 := t1.NumField()
	if n0 != n1 {
		return false
	}

	for i := range n0 {
		f0 := t0.Field(i)
		f1 := t1.Field(i)
		if !isTypeEqual(f0.Type, f1.Type) {
			return false
		}
	}
	return true
}

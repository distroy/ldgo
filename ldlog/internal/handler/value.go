/*
 * Copyright (C) distroy
 */

package handler

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"
)

func init() {
	checkTypeEqual(reflect.TypeOf(Value{}), reflect.TypeOf(slog.Value{}))
}

func getValuePtr(v *slog.Value) *Value { return toType[*Value](v) }
func getValue(v slog.Value) Value      { return *getValuePtr(&v) }

func countEmptyGroups(as []Attr) int {
	n := 0
	for _, a := range as {
		if a.Value.IsEmptyGroup() {
			n++
		}
	}
	return n
}

func BoolValue(v bool) Value         { return getValue(slog.BoolValue(v)) }
func StringValue(value string) Value { return getValue(slog.StringValue(value)) }
func Float64Value(v float64) Value   { return getValue(slog.Float64Value(v)) }

func IntValue(v int) Value       { return Int64Value(int64(v)) }
func Int64Value(v int64) Value   { return getValue(slog.Int64Value(v)) }
func Uint64Value(v uint64) Value { return getValue(slog.Uint64Value(v)) }

func TimeValue(v time.Time) Value         { return getValue(slog.TimeValue(v)) }
func DurationValue(v time.Duration) Value { return getValue(slog.DurationValue(v)) }

func GroupValue(as ...Attr) Value { return getValue(slog.GroupValue(GetSAttrs(as)...)) }
func AnyValue(v any) Value        { return getValue(slog.AnyValue(v)) }

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

func (v *Value) Get() *slog.Value { return toType[*slog.Value](v) }

func (v *Value) Kind() Kind { return v.Get().Kind() }

func (v *Value) Any() any       { return v.Get().Any() }
func (v *Value) String() string { return v.Get().String() }

func (v *Value) Int64() int64     { return v.Get().Int64() }
func (v *Value) Uint64() uint64   { return v.Get().Uint64() }
func (v *Value) Float64() float64 { return v.Get().Float64() }

func (v *Value) Bool() bool              { return v.Get().Bool() }
func (v *Value) Duration() time.Duration { return v.Get().Duration() }
func (v *Value) Time() time.Time         { return v.Get().Time() }

func (v *Value) LogValuer() slog.LogValuer { return v.Get().LogValuer() }
func (v *Value) Group() []Attr             { return GetAttrs(v.Get().Group()) }

func (v *Value) Equal(w Value) bool { return v.Get().Equal(*w.Get()) }

func (v *Value) Resolve() (rv Value) { return toType[Value](v.Get().Resolve()) }

// IsEmptyGroup reports whether v is a group that has no attributes.
func (v *Value) IsEmptyGroup() bool {
	if v.Kind() != KindGroup {
		return false
	}
	// We do not need to recursively examine the group's Attrs for emptiness,
	// because GroupValue removed them when the group was constructed, and
	// groups are immutable.
	return len(v.Group()) == 0
}

func (v *Value) Append(dst []byte) []byte {
	switch v.Kind() {
	case KindString:
		return append(dst, v.String()...)
	case KindInt64:
		return strconv.AppendInt(dst, int64(v.num), 10)
	case KindUint64:
		return strconv.AppendUint(dst, v.num, 10)
	case KindFloat64:
		return strconv.AppendFloat(dst, v.Float64(), 'g', -1, 64)
	case KindBool:
		return strconv.AppendBool(dst, v.Bool())
	case KindDuration:
		return append(dst, v.Duration().String()...)
	case KindTime:
		return append(dst, v.Time().String()...)
	case KindGroup:
		return fmt.Append(dst, v.Group())
	case KindAny, KindLogValuer:
		return fmt.Append(dst, v.any)
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
}

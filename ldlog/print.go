/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"

	"github.com/distroy/ldgo/v3/ldlog/internal/_buffer"
	"github.com/distroy/ldgo/v3/ldref"
)

var (
	_ fmt.Stringer   = (*printWrapper)(nil)
	_ slog.LogValuer = (*printWrapper)(nil)
)

type printWrapper struct {
	args []any
}

func (w printWrapper) String() string {
	return sprintln(w.args)
}

func (w printWrapper) LogValue() Value {
	return slog.StringValue(w.String())
}

func pw(args []any) printWrapper { return printWrapper{args: args} }

func sprintln(args []any) string {
	if len(args) == 0 {
		return ""
	}

	buf := _buffer.NewBuffer()

	fprintArg(buf, args[0])
	for _, arg := range args[1:] {
		buf.WriteByte(' ')
		fprintArg(buf, arg)
	}

	buf.TrimNewline()
	text := buf.String()
	buf.Free()

	return text
}

func fprintArg(b *_buffer.Buffer, val any) {
	switch v := val.(type) {
	case fmt.Stringer:
		b.WriteString(v.String())
		return

	case error:
		b.WriteString(v.Error())
		return
	}

	ref := reflect.ValueOf(val)
	if ref.Kind() == reflect.Ptr {
		if ref.Pointer() == 0 {
			fprintPointer(b, ref)
			return
		}
		ref = ref.Elem()
	}

	switch ref.Kind() {
	case reflect.Array, reflect.Slice:
		fprintSlice(b, ref)

	case reflect.Map:
		fprintMap(b, ref)

	case reflect.Struct:
		fprintStruct(b, ref)

	case reflect.String:
		b.WriteString(ref.String())

	case reflect.Bool:
		b.AppendBool(ref.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.AppendInt(ref.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.AppendUint(ref.Uint())

	case reflect.Float64:
		b.AppendFloat(ref.Float(), 64)

	case reflect.Float32:
		b.AppendFloat(ref.Float(), 32)

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		fprintPointer(b, ref)

	default:
		fmt.Fprint(b, ref.Interface())
	}
}

func fprintSlice(b *_buffer.Buffer, v reflect.Value) {
	b.WriteString("[")
	for i := 0; i < v.Len(); i++ {
		if i != 0 {
			b.WriteString(", ")
		}
		fprintArg(b, v.Index(i).Interface())
	}
	b.WriteString("]")
}

func fprintPointer(b *_buffer.Buffer, v reflect.Value) {
	p := v.Pointer()

	b.WriteByte('(')
	b.WriteString(v.Type().String())
	b.WriteString(")(")
	if p == 0 {
		b.WriteString("nil")
	} else {
		fmt.Fprintf(b, "0x%x", p)
	}
	b.WriteByte(')')
}

func fprintStruct(b *_buffer.Buffer, v reflect.Value) {
	b.WriteByte('{')
	for i := 0; i < v.NumField(); i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		if name := v.Type().Field(i).Name; name != "" {
			b.WriteString(name)
			b.WriteByte(':')
		}
		field := v.Field(i)
		fprintArg(b, field.Interface())
	}
	b.WriteByte('}')
}

func fprintMap(b *_buffer.Buffer, val reflect.Value) {
	m := make([][2]reflect.Value, 0, val.Len())
	for it := val.MapRange(); it.Next(); {
		m = append(m, [2]reflect.Value{it.Key(), it.Value()})
	}

	sort.Sort(sortedMap(m))

	b.WriteString("map[")
	for i, kv := range m {
		if i > 0 {
			b.WriteByte(',')
		}
		fprintArg(b, kv[0].Interface())
		b.WriteByte(':')
		fprintArg(b, kv[1].Interface())
	}
	b.WriteByte(']')
}

type sortedMap [][2]reflect.Value

func (o sortedMap) Len() int           { return len(o) }
func (o sortedMap) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o sortedMap) Less(i, j int) bool { return ldref.CompareReflect(o[i][0], o[j][0]) <= 0 }

/*
 * Copyright (C) distroy
 */

package attr

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"github.com/distroy/ldgo/v3/internal/jsontag"
	"github.com/distroy/ldgo/v3/ldcmp"
	"github.com/distroy/ldgo/v3/ldlog/internal/buffer"
	"github.com/distroy/ldgo/v3/ldsort"
)

func BriefReflect(key string, val any) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Any(key, brief_reflect_t{val})
}

func mapkey2str(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	case reflect.String:
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Complex64:
		return strconv.FormatComplex(v.Complex(), 'f', -1, 64)
	case reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'f', -1, 128)
	}

	return "<unknown>"
}

func addBriefRef(b *buffer.Buffer, v reflect.Value) {
	switch vv := v.Interface().(type) {
	case time.Time:
		// b.AppendTime(vv, time.RFC3339)
		b.AppendTime(vv, "")
		return

	case time.Duration:
		b.AppendString(quote(vv.String()))
		return

	case fmt.Stringer:
		b.AppendString(quote(vv.String()))
		return

	case error:
		b.AppendString(quote(vv.Error()))
		return
	}

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			b.AppendString(nil_t{}.String())
			return
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Invalid:
		b.AppendString(nil_t{}.String())
		return

	case reflect.String:
		addBriefStr(b, v.String())
		return

	case reflect.Bool:
		b.AppendBool(v.Bool())
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.AppendInt(v.Int())
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b.AppendUint(v.Uint())
		return

	case reflect.Float32:
		b.AppendFloat(v.Float(), 32)
		return
	case reflect.Float64:
		b.AppendFloat(v.Float(), 64)
		return

	case reflect.Complex64:
		b.AppendComplex(v.Complex(), 64)
		return
	case reflect.Complex128:
		b.AppendComplex(v.Complex(), 128)
		return

	case reflect.Slice:
		switch vv := v.Interface().(type) {
		case []byte:
			addBriefStr(b, b2s(vv))
			return
		}
		fallthrough
	case reflect.Array:
		addBriefRefSlice(b, v)
		return

	case reflect.Struct:
		addBriefRefStruct(b, v)
		return

	case reflect.Map:
		addBriefRefMap(b, v)
		return
	}
}

func addBriefRefSlice(b *buffer.Buffer, v reflect.Value) {
	n := briefArrayLen
	l := v.Len()
	if l <= n {
		b.AppendByte('[')
		addBriefRefSliceData(b, v, l)
		b.AppendByte(']')
		return
	}

	b.AppendByte('{')

	b.AppendString(tagLen)
	b.AppendByte(':')
	b.AppendInt(int64(l))
	b.AppendByte(',')

	b.AppendString(tagType)
	b.AppendByte(':')
	b.AppendString("array")
	b.AppendByte(',')

	b.AppendString(tagBrief)
	b.AppendByte(':')

	b.AppendByte('[')
	addBriefRefSliceData(b, v, n)
	b.AppendByte(']')

	b.AppendByte('}')
}

func addBriefRefSliceData(b *buffer.Buffer, v reflect.Value, l int) {
	for i := range l {
		if i > 0 {
			b.AppendByte(',')
		}
		f := v.Index(i)
		addBriefRef(b, f)
	}
}

func addBriefRefStruct(b *buffer.Buffer, v reflect.Value) {
	typ := v.Type()
	s := jsontag.Get(typ)
	b.AppendByte('{')
	addBriefRefStructData(b, v, s)
	b.AppendByte('}')
}

func addBriefRefStructData(b *buffer.Buffer, v reflect.Value, s *jsontag.Struct) {
	for i := range s.NumField() {
		ft := s.Field(i)
		fv := v.Field(i)
		addBriefRefStructField(b, fv, ft)
	}
}

func addBriefRefStructField(b *buffer.Buffer, v reflect.Value, f *jsontag.Field) {
	if f.Field.Anonymous {
		addBriefRefStructFieldEmbeded(b, v, f)
		return
	}
	if f.OmitEmpty && v.IsZero() {
		return
	}
	addBriefRefKeyValue(b, f.Name, v)
}

func addBriefRefStructFieldEmbeded(b *buffer.Buffer, v reflect.Value, f *jsontag.Field) {
	if f.Type.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	typ := v.Type()
	s := jsontag.Get(typ)
	addBriefRefStructData(b, v, s)
}

func addBriefRefMap(b *buffer.Buffer, v reflect.Value) {
	n := briefMapLen
	l := v.Len()
	if n == 0 || l <= n {
		b.AppendByte('{')
		addBriefRefMapData(b, v, l)
		b.AppendByte('}')
		return
	}

	b.AppendByte('{')

	b.AppendString(tagLen)
	b.AppendByte(':')
	b.AppendInt(int64(l))
	b.AppendByte(',')

	b.AppendString(tagType)
	b.AppendByte(':')
	b.AppendString("map")
	b.AppendByte(',')

	b.AppendString(tagBrief)
	b.AppendByte(':')

	b.AppendByte('{')
	addBriefRefMapData(b, v, n)
	b.AppendByte('}')

	b.AppendByte('}')
}

func addBriefRefMapData(b *buffer.Buffer, v reflect.Value, l int) {
	type data struct {
		Key string
		Val reflect.Value
	}
	s := make([]data, 0, v.Len())
	for it := v.MapRange(); it.Next(); {
		k := mapkey2str(it.Key())
		v := it.Value()
		s = append(s, data{Key: k, Val: v})
	}

	ldsort.Sort(s, func(a, b data) int { return ldcmp.CompareString(a.Key, b.Key) })
	for i := range l {
		d := &s[i]
		addBriefRefKeyValue(b, d.Key, d.Val)
	}
}

func addBriefRefKeyValue(b *buffer.Buffer, key string, val reflect.Value) {
	l := b.Len()
	if l > 0 {
		switch b.Bytes()[l-1] {
		case '{':
			b.AppendByte(',')
		}
	}
	b.AppendString(key)
	b.AppendByte(':')
	addBriefRef(b, val)
}

type brief_reflect_t struct {
	val any
}

func (p brief_reflect_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_reflect_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_reflect_t) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (p brief_reflect_t) String() string {
	if p.val == nil {
		return nil_t{}.String()
	}

	buf := getBuf()
	defer buf.Free()

	v, ok := p.val.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(v)
	}
	addBriefRef(buf, v)
	return buf.String()
}

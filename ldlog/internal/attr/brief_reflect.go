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
		b.AppendTime(vv, time.RFC3339)
		return

	case time.Duration:
		b.AppendString(quote(vv.String()))
		return

	case fmt.Stringer:
		b.AppendString(quote(vv.String()))
		return

	case error:
		b.AppendString(quote(vv.Error()))
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
		b.AppendString(v.String())
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
		return enc.AppendObject(&briefReflectStruct{Val: v})

	case reflect.Map:
		return enc.AppendObject(&briefReflectMap{Val: v})
	}
}

func addBriefRefSlice(b *buffer.Buffer, v reflect.Value) {
	n := briefArrayLen
	l := v.Len()
	if l > n {
		b.AppendByte('{')

		b.AppendString(tagLen)
		b.AppendByte(':')
		b.AppendInt(int64(l))
		b.AppendByte(',')

		b.AppendString(tagType)
		b.AppendByte(':')
		b.AppendString("array")
		b.AppendByte(',')
	}

	n = min(n, l)
	for i := range n {
		if i > 0 {
			b.AppendByte(',')
		}

		f := v.Index(i)
		addBriefRef(b, f)
	}

	if l > n {
		b.AppendByte('}')
	}
}

func addBriefRefStruct(b *buffer.Buffer, v reflect.Value) {
	typ := v.Type()
	s := jsontag.Get(typ)
	b.AppendByte('{')
	first := true
	for i := range s.NumField() {
		if !first {
			b.AppendByte(',')
		}
		fs := s.Field(i)
		if fs.Field.Anonymous {
			continue
		}
	}
}

type brief_reflect_t struct {
	val any
}

func (p brief_reflect_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_reflect_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_reflect_t) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (p brief_reflect_t) String() string {
	buf := getBuf()
	defer buf.Free()
	if p.val == nil {
		buf.WriteString(nil_t{}.String())
		return buf.String()
	}
}

type briefReflectMap struct {
	Val reflect.Value
	Len int
}

func (p *briefReflectMap) MarshalLogObject(enc ObjectEncoder) error {
	if p.Len > 0 {
		return p.marshalLogObject(enc, p.Len)
	}

	n := briefMapLen
	l := p.Val.Len()
	if l <= n {
		return p.marshalLogObject(enc, l)
	}

	enc.AddInt(tagLen, l)
	enc.AddString(tagType, "map")
	return enc.AddObject(tagBrief, &briefReflectMap{Val: p.Val, Len: n})
}

func (p *briefReflectMap) marshalLogObject(enc ObjectEncoder, n int) error {
	type data struct {
		Key string
		Val reflect.Value
	}
	l := make([]data, 0, p.Val.Len())
	for it := p.Val.MapRange(); it.Next(); {
		k := mapkey2str(it.Key())
		v := it.Value()
		l = append(l, data{Key: k, Val: v})
	}

	ldsort.Sort(l, func(a, b data) int { return ldcmp.CompareString(a.Key, b.Key) })
	for i := 0; i < n; i++ {
		d := &l[i]
		err := AddRef2Log(enc, d.Key, d.Val)
		if err != nil {
			return err
		}
	}
	return nil
}

type briefReflectArray struct {
	Val reflect.Value
	Len int
}

func (p *briefReflectArray) MarshalLogObject(enc ObjectEncoder) error {
	enc.AddInt(tagLen, p.Val.Len())
	enc.AddString(tagType, "array")
	return enc.AddArray(tagBrief, p)
}

func (p *briefReflectArray) MarshalLogArray(enc ArrayEncoder) error {
	for i := 0; i < p.Len; i++ {
		v := p.Val.Index(i)
		err := AppendRef2Log(enc, v)
		if err != nil {
			return err
		}
	}
	return nil
}

type briefReflectStruct struct {
	Val reflect.Value
}

func (p *briefReflectStruct) MarshalLogObject(enc ObjectEncoder) error {
	return marshalReflectStruct(enc, p.Val)
}

func marshalReflectStruct(enc ObjectEncoder, obj reflect.Value) error {
	typ := obj.Type()
	s := jsontag.Get(typ)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		err := marshalReflectStructField(enc, obj, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func marshalReflectStructField(enc ObjectEncoder, obj reflect.Value, f *jsontag.Field) error {
	k := f.Name
	v := obj.Field(f.Index)

	// log.Printf(" === field begin. field: %s", k)
	// defer log.Printf(" === field end. field: %s", k)

	// if !f.Field.IsExported() {
	// 	addr := unsafe.Pointer(v.UnsafeAddr())
	// 	v = reflect.NewAt(v.Type(), addr).Elem()
	// }

	if f.Field.Anonymous {
		return marshalReflectStructEmbedded(enc, v, f)
	}

	if f.OmitEmpty && v.IsZero() {
		return nil
	}

	return AddRef2Log(enc, k, v)
}

func marshalReflectStructEmbedded(enc ObjectEncoder, v reflect.Value, f *jsontag.Field) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	// if !f.Field.IsExported() {
	// 	addr := unsafe.Pointer(v.UnsafeAddr())
	// 	v = reflect.NewAt(v.Type(), addr).Elem()
	// }
	return marshalReflectStruct(enc, v)
}

/*
 * Copyright (C) distroy
 */

package field

import (
	"reflect"
	"strconv"
	"time"

	"github.com/distroy/ldgo/v2/internal/jsontag"
	"github.com/distroy/ldgo/v2/ldcmp"
	"github.com/distroy/ldgo/v2/ldsort"
	"go.uber.org/zap"
)

func BriefReflect(key string, val interface{}) Field {
	if val == nil {
		return zap.Reflect(key, val)
	}
	return zap.Inline(&BriefReflectType{Key: key, Val: val})
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

var (
	_ ArrayEncoder = (*objectEncoder)(nil)
)

type objectEncoder struct {
	Key string
	Enc ObjectEncoder
}

func (p *objectEncoder) AppendBool(v bool) { p.Enc.AddBool(p.Key, v) }

func (p *objectEncoder) AppendString(v string)     { p.Enc.AddString(p.Key, v) }
func (p *objectEncoder) AppendByteString(v []byte) { p.Enc.AddByteString(p.Key, v) }

func (p *objectEncoder) AppendFloat32(v float32) { p.Enc.AddFloat32(p.Key, v) }
func (p *objectEncoder) AppendFloat64(v float64) { p.Enc.AddFloat64(p.Key, v) }

func (p *objectEncoder) AppendComplex64(v complex64)   { p.Enc.AddComplex64(p.Key, v) }
func (p *objectEncoder) AppendComplex128(v complex128) { p.Enc.AddComplex128(p.Key, v) }

func (p *objectEncoder) AppendInt(v int)     { p.Enc.AddInt(p.Key, v) }
func (p *objectEncoder) AppendInt8(v int8)   { p.Enc.AddInt8(p.Key, v) }
func (p *objectEncoder) AppendInt16(v int16) { p.Enc.AddInt16(p.Key, v) }
func (p *objectEncoder) AppendInt32(v int32) { p.Enc.AddInt32(p.Key, v) }
func (p *objectEncoder) AppendInt64(v int64) { p.Enc.AddInt64(p.Key, v) }

func (p *objectEncoder) AppendUint(v uint)       { p.Enc.AddUint(p.Key, v) }
func (p *objectEncoder) AppendUint8(v uint8)     { p.Enc.AddUint8(p.Key, v) }
func (p *objectEncoder) AppendUint16(v uint16)   { p.Enc.AddUint16(p.Key, v) }
func (p *objectEncoder) AppendUint32(v uint32)   { p.Enc.AddUint32(p.Key, v) }
func (p *objectEncoder) AppendUint64(v uint64)   { p.Enc.AddUint64(p.Key, v) }
func (p *objectEncoder) AppendUintptr(v uintptr) { p.Enc.AddUintptr(p.Key, v) }

func (p *objectEncoder) AppendTime(v time.Time)         { p.Enc.AddTime(p.Key, v) }
func (p *objectEncoder) AppendDuration(v time.Duration) { p.Enc.AddDuration(p.Key, v) }

func (p *objectEncoder) AppendArray(v ArrayMarshaler) error   { return p.Enc.AddArray(p.Key, v) }
func (p *objectEncoder) AppendObject(v ObjectMarshaler) error { return p.Enc.AddObject(p.Key, v) }

func (p *objectEncoder) AppendReflected(v interface{}) error { return p.Enc.AddReflected(p.Key, v) }

func AddRef2Log(enc ObjectEncoder, k string, v reflect.Value) error {
	oe := &objectEncoder{
		Key: k,
		Enc: enc,
	}
	return AppendRef2Log(oe, v)
}

func AppendRef2Log(enc ArrayEncoder, v reflect.Value) error {
	switch vv := v.Interface().(type) {
	case time.Time:
		enc.AppendTime(vv)
		return nil

	case time.Duration:
		enc.AppendDuration(vv)
		return nil

	// case fmt.Stringer:
	// 	AppendStr2Log(enc, vv.String())
	// 	return nil

	case error:
		enc.AppendString(vv.Error())
	}

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			// enc.AppendString(tagNil)
			enc.AppendReflected(nil)
			return nil
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Invalid:
		// enc.AppendString(tagNil)
		enc.AppendReflected(nil)
		return nil

	case reflect.String:
		return AppendStr2Log(enc, v.String())

	case reflect.Bool:
		enc.AppendBool(v.Bool())
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.AppendInt64(v.Int())
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		enc.AppendUint64(v.Uint())
		return nil

	case reflect.Float32, reflect.Float64:
		enc.AppendFloat64(v.Float())
		return nil

	case reflect.Complex64, reflect.Complex128:
		enc.AppendComplex128(v.Complex())
		return nil

	case reflect.Slice:
		switch vv := v.Interface().(type) {
		case []byte:
			return AppendStr2Log(enc, bytes2str(vv))
		}
		fallthrough
	case reflect.Array:
		n := briefArrayLen
		l := v.Len()
		if l <= n {
			return enc.AppendArray(&briefReflectArray{Val: v, Len: l})
		}
		return enc.AppendObject(&briefReflectArray{Val: v, Len: n})

	case reflect.Struct:
		return enc.AppendObject(&briefReflectStruct{Val: v})

	case reflect.Map:
		return enc.AppendObject(&briefReflectMap{Val: v})
	}
	return nil
}

type BriefReflectType struct {
	Key string
	Val interface{}
}

func (p *BriefReflectType) MarshalLogObject(enc ObjectEncoder) error {
	vv := p.Val
	v := reflect.ValueOf(vv)
	return AddRef2Log(enc, p.Key, v)
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

/*
 * Copyright (C) distroy
 */

package _handler

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/distroy/ldgo/v3/ldconv"
)

type LogTextAppender interface {
	AppendLogText(b []byte) ([]byte, error)
}

// func s2b(b string) []byte { return ldconv.StrToBytesUnsafe(b) }
func b2s(b []byte) string   { return ldconv.BytesToStrUnsafe(b) }
func quote(s string) string { return strconv.Quote(s) }

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

func appendTextValue(s *handleState, v Value) error {
	switch v.Kind() {
	case KindString:
		s.appendString(v.String())
	case KindTime:
		s.appendTime(v.Time())
	case KindAny:
		vv := v.DirectlyAny()
		switch m := vv.(type) {
		case io.WriterTo:
			_, err := m.WriteTo(s.buf)
			return err

		case encoding.TextMarshaler:
			data, err := m.MarshalText()
			if err != nil {
				return err
			}
			// TODO: avoid the conversion to string.
			s.appendString(b2s(data))
			return nil
		}
		if bs, ok := byteSlice(vv); ok {
			// As of Go 1.19, this only allocates for strings longer than 32 bytes.
			s.buf.WriteString(quote(b2s(bs)))
			return nil
		}
		s.appendString(fmt.Sprintf("%+v", vv))
		// s.appendStringWithoutQuote(fmt.Sprintf("%+v", v.any))
	default:
		*s.buf = v.Append(*s.buf)
	}
	return nil
}

func byteSlice(a any) ([]byte, bool) {
	if bs, ok := a.([]byte); ok {
		return bs, true
	}
	// Like Printf's %s, we allow both the slice type and the byte element type to be named.
	t := reflect.TypeOf(a)
	if t != nil && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(a).Bytes(), true
	}
	return nil, false
}

// //go:linkname appendJSONValue log/slog.appendJSONValue
// func appendJSONValue(s *handleState, v Value) error

func appendJSONValue(s *handleState, v Value) error {
	switch v.Kind() {
	case KindString:
		s.appendString(v.String())
	case KindInt64:
		*s.buf = strconv.AppendInt(*s.buf, v.Int64(), 10)
	case KindUint64:
		*s.buf = strconv.AppendUint(*s.buf, v.Uint64(), 10)
	case KindFloat64:
		// json.Marshal is funny about floats; it doesn't
		// always match strconv.AppendFloat. So just call it.
		// That's expensive, but floats are rare.
		if err := appendJSONMarshal(s.buf, v.Float64()); err != nil {
			return err
		}
	case KindBool:
		*s.buf = strconv.AppendBool(*s.buf, v.Bool())
	case KindDuration:
		// Do what json.Marshal does.
		*s.buf = strconv.AppendInt(*s.buf, int64(v.Duration()), 10)
	case KindTime:
		s.appendTime(v.Time())
	case KindAny:
		a := v.Any()
		_, jm := a.(json.Marshaler)
		if err, ok := a.(error); ok && !jm {
			s.appendString(err.Error())
		} else {
			return appendJSONMarshal(s.buf, a)
		}
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
	return nil
}

func appendJSONMarshal(buf *Buffer, v any) error {
	// Use a json.Encoder to avoid escaping HTML.
	// var bb bytes.Buffer
	bb := newBuffer()
	defer func() {
		bb.Free()
	}()
	enc := json.NewEncoder(bb)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	bs := bb.Bytes()
	buf.Write(bs[:len(bs)-1]) // remove final newline
	return nil
}

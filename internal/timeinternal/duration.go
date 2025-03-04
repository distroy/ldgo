/*
 * Copyright (C) distroy
 */

package timeinternal

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

func DurationMarshalJson(d time.Duration) ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+4)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

func DurationUnmarshalJson(b []byte) (time.Duration, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("unexpected end of JSON input")
	}

	switch b[0] {
	case '"', '\'', '`':
		return durationUnmarshalJsonByStr(b)
	}

	return durationUnmarshalJsonByInt(b)
}

func durationUnmarshalJSONError(b []byte) error {
	return fmt.Errorf("invalid duration: %s", b)
}

func durationUnmarshalJsonByStr(b []byte) (time.Duration, error) {
	str := unsafe.String(unsafe.SliceData(b), len(b))
	vv, err := strconv.Unquote(str)
	if err != nil {
		return 0, durationUnmarshalJSONError(b)
	}
	dur, err := time.ParseDuration(vv)
	if err != nil {
		return 0, durationUnmarshalJSONError(b)
	}
	return dur, nil
}

func durationUnmarshalJsonByInt(b []byte) (time.Duration, error) {
	if b[0] != '0' {
		str := unsafe.String(unsafe.SliceData(b), len(b))
		if i64, err := strconv.ParseInt(str, 10, 64); err == nil {
			return time.Duration(i64), nil
		}
		return 0, durationUnmarshalJSONError(b)
	}

	if len(b) == 1 {
		return 0, nil
	}

	switch b[1] {
	case 'x', 'X':
		bb := b[2:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 16, 64); err == nil {
			return time.Duration(i64), nil
		}

	case 'o', 'O':
		bb := b[2:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 8, 64); err == nil {
			return time.Duration(i64), nil
		}

	default:
		bb := b[1:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 8, 64); err == nil {
			return time.Duration(i64), nil
		}
	}

	return 0, durationUnmarshalJSONError(b)
}

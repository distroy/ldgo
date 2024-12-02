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

func DurationMarshalJSON(d time.Duration) ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+4)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

func DurationUnmarshalJSON(b []byte) (time.Duration, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("unexpected end of JSON input")
	}

	str := unsafe.String(unsafe.SliceData(b), len(b))
	switch b[0] {
	case '"', '\'', '`':
		vv, err := strconv.Unquote(str)
		if err != nil {
			return 0, fmt.Errorf("invalid duration: %s", b)
		}
		dur, err := time.ParseDuration(vv)
		if err != nil {
			return 0, fmt.Errorf("invalid duration: %s", b)
		}
		return dur, nil

	}

	if i64, err := strconv.ParseInt(str, 10, 64); err == nil {
		return time.Duration(i64), nil
	}
	// if f64, err := strconv.ParseFloat(str, 64); err != nil {
	// 	*d = Duration(f64)
	// 	return nil
	// }
	return 0, fmt.Errorf("invalid duration: %s", b)
}

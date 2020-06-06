/*
 * Copyright (C) distroy
 */

package ldgo

import (
	"strconv"
)

func ToByte(v interface{}) (byte, error) { return ToUint8(v) }

func ToBool(_v interface{}) (bool, error) {
	switch v := _v.(type) {
	case bool:
		return v, nil

	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil

	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil

	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil

	case []byte:
		return convBool(v)
	case string:
		return convBool(StrToBytesUnsafe(v))
	}
	return false, _ERR_UNKOWN_TYPE
}

func ToInt(v interface{}) (int, error) {
	n, err := ToInt64(v)
	return int(n), err
}
func ToInt8(v interface{}) (int8, error) {
	n, err := ToInt64(v)
	return int8(n), err
}
func ToInt16(v interface{}) (int16, error) {
	n, err := ToInt64(v)
	return int16(n), err
}
func ToInt32(v interface{}) (int32, error) {
	n, err := ToInt64(v)
	return int32(n), err
}

func ToInt64(_v interface{}) (int64, error) {
	switch v := _v.(type) {
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil

	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil

	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil

	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil

	case []byte:
		return convInt(v)
	case string:
		return convInt(StrToBytesUnsafe(v))
	}
	return 0, _ERR_UNKOWN_TYPE
}

func ToUint(v interface{}) (uint, error) {
	n, err := ToUint64(v)
	return uint(n), err
}
func ToUint8(v interface{}) (uint8, error) {
	n, err := ToUint64(v)
	return uint8(n), err
}
func ToUint16(v interface{}) (uint16, error) {
	n, err := ToUint64(v)
	return uint16(n), err
}
func ToUint32(v interface{}) (uint32, error) {
	n, err := ToUint64(v)
	return uint32(n), err
}

func ToUint64(_v interface{}) (uint64, error) {
	switch v := _v.(type) {
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil

	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil

	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil

	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil

	case []byte:
		return convUint(v)
	case string:
		return convUint(StrToBytesUnsafe(v))
	}
	return 0, _ERR_UNKOWN_TYPE
}

func ToFloat32(v interface{}) (float32, error) {
	n, err := ToFloat64(v)
	return float32(n), err
}

func ToFloat64(_v interface{}) (float64, error) {
	switch v := _v.(type) {
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil

	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil

	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil

	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil

	case []byte:
		return convFloat(v)
	case string:
		return convFloat(StrToBytesUnsafe(v))
	}
	return 0, _ERR_UNKOWN_TYPE
}

func ToString(_v interface{}) (string, error) {
	switch v := _v.(type) {
	case bool:
		return strconv.FormatBool(v), nil

	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(int64(v), 10), nil

	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(v), 10), nil

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil

	case []byte:
		return BytesToStrUnsafe(v), nil
	case string:
		return v, nil
	}
	return "", _ERR_UNKOWN_TYPE
}

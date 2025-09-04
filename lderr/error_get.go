/*
 * Copyright (C) distroy
 */

package lderr

import (
	"context"
	"net/http"
)

func GetCode(err error, def ...int) int {
	if err == nil {
		return 0
	}
	if code, ok := getCode(err); ok {
		return code
	}
	if len(def) > 0 {
		return def[0]
	}
	return errCodeUnkown
}
func getCode(err error) (int, bool) {
	return unwrap(err, func(err error) (int, bool) {
		if e := getMatchError(err); e != nil {
			return e.Code(), true
		}
		if x, ok := err.(interface{ Code() int }); ok {
			return x.Code(), true
		}
		return 0, false
	})
}

func unwrap[T any](err error, get func(err error) (T, bool)) (T, bool) {
	zero := func() (x T) { return }
	for {
		if r, ok := get(err); ok {
			return r, true
		}

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err == nil {
				return zero(), false
			}
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				r, ok := get(err)
				if ok {
					return r, true
				}
			}
			return zero(), false
		default:
			return zero(), false
		}
	}
}

func GetStatus(err error, def ...int) int {
	if err == nil {
		return http.StatusOK
	}
	if status, ok := getStatus(err); ok {
		return status
	}
	if len(def) > 0 {
		return def[0]
	}
	return errStatusUnkonw
}
func getStatus(err error) (int, bool) {
	return unwrap(err, func(err error) (int, bool) {
		if e := getMatchError(err); e != nil {
			return e.Status(), true
		}
		if x, ok := err.(interface{ Status() int }); ok {
			return x.Status(), true
		}
		return 0, false
	})
}

func GetMessage(err error) string {
	if err == nil {
		return errMessageSucess
	}
	return err.Error()
}

func GetDetails(err error) []string {
	if err == nil {
		return nil
	}
	details, _ := getDetails(err)
	return details
}
func getDetails(err error) ([]string, bool) {
	return unwrap(err, func(err error) ([]string, bool) {
		switch v := err.(type) {
		case interface{ Details() []string }:
			return v.Details(), true

		case interface{ Detail() string }:
			return []string{v.Detail()}, true
		}
		return nil, false
	})
}

func getMatchError(err error) Error {
	switch err {
	case nil:
		return ErrSuccess

	case context.Canceled:
		return ErrCtxCanceled

	case context.DeadlineExceeded:
		return ErrCtxDeadlineExceeded
	}
	return nil
}

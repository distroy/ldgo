/*
 * Copyright (C) distroy
 */

package lderr

import (
	"errors"
)

type Error interface {
	error

	Status() int
	Code() int
	Details() []string
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, Is examines err followed by a
// depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err, target error) bool {
	if err == nil {
		return IsSuccess(target)
	}
	if target == nil {
		return IsSuccess(err)
	}
	return errors.Is(err, target)
}

func In(err error, targets ...error) bool {
	for _, target := range targets {
		if Is(err, target) {
			return true
		}
	}
	return false
}

func IsSuccess(err error) bool {
	if err == nil {
		return true
	}
	if GetCode(err) == 0 {
		return true
	}
	return false
}

func New(status, code int, message string) Error {
	return commError{
		error:  strError(message),
		status: status,
		code:   code,
	}
}

func Wrap(err error, def ...Error) Error {
	if err == nil {
		return nil
	}

	if v, ok := err.(Error); ok {
		return v
	}

	if e := getMatchError(err); e != nil {
		return e
	}

	d := ErrUnkown
	if len(def) != 0 {
		d = def[0]
	}

	e := commError{
		error:  err,
		status: d.Status(),
		code:   d.Code(),
	}
	return newWithDetails(e, GetDetails(err))
}

func Override(err error, message string) Error {
	e := commError{
		error:  strError(message),
		status: GetStatus(err),
		code:   GetCode(err),
	}
	return newWithDetails(e, GetDetails(err))
}

func newWithDetails(err commError, details []string) Error {
	if len(details) == 0 {
		return err
	}
	return &detailsError{
		commError: err,
		details:   details,
	}
}

type commError struct {
	error

	status int
	code   int
}

func (e commError) Status() int       { return e.status }
func (e commError) Code() int         { return e.code }
func (e commError) Unwrap() error     { return e.error }
func (e commError) Details() []string { return nil }
func (e commError) Is(target error) bool {
	if err, _ := target.(interface{ Code() int }); err != nil && e.Code() == err.Code() {
		return true
	}
	// return Is(e.error, target)
	return false
}

type strError string

func (e strError) Error() string { return string(e) }

func WithDetail(err error, details ...string) Error {
	return WithDetails(err, details)
}

func WithDetails(err error, details []string) Error {
	t := GetDetails(err)

	if len(details)+len(t) == 0 {
		return Wrap(err)
	}

	if len(details) == 0 {
		return &detailsError{
			commError: commError{
				error:  err,
				status: GetStatus(err),
				code:   GetCode(err),
			},
			details: t,
		}
	}

	if len(t) == 0 {
		return &detailsError{
			commError: commError{
				error:  err,
				status: GetStatus(err),
				code:   GetCode(err),
			},
			details: details,
		}
	}

	d := make([]string, 0, len(details)+len(t))
	d = append(d, t...)
	d = append(d, details...)

	return &detailsError{
		commError: commError{
			error:  err,
			status: GetStatus(err),
			code:   GetCode(err),
		},
		details: d,
	}
}

type detailsError struct {
	commError

	details []string
}

func (e *detailsError) Details() []string { return e.details }
func (e *detailsError) Unwrap() error     { return e.commError }
func (e *detailsError) Is(target error) bool {
	if err, _ := target.(interface{ Code() int }); err != nil && e.Code() == err.Code() {
		return true
	}
	// return Is(e.commError, target)
	return false
}

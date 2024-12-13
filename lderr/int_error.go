/*
 * Copyright (C) distroy
 */

package lderr

var intErrors = make(map[intError]commError)

func newInt(status, code int, message string) Error {
	return newIntByErr(status, code, strError(message))
}

func newIntByErr(status, code int, message error) Error {
	err := intError(code)
	if _, ok := intErrors[err]; !ok {
		intErrors[err] = commError{
			error:  message,
			status: status,
			code:   code,
		}
		return err
	}
	return err
}

type intError int

func (e intError) err() Error {
	err, ok := intErrors[e]
	if ok {
		return err
	}
	return commError{
		error:  strError(errMessageUnkonw),
		status: errStatusUnkonw,
		code:   e.Code(),
	}
}

func (e intError) Code() int         { return int(e) }
func (e intError) Details() []string { return nil }
func (e intError) Error() string     { return e.err().Error() }
func (e intError) Status() int       { return e.err().Status() }
func (e intError) Unwrap() error     { return e.err() }
func (e intError) Is(target error) bool {
	if err, _ := target.(interface{ Code() int }); err != nil && e.Code() == err.Code() {
		return true
	}
	// return Is(e.Unwrap(), target)
	return false
}

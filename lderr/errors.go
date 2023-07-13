/*
 * Copyright (C) distroy
 */

package lderr

import "net/http"

// error definitions: [-1, -999]
var (
	ErrSuccess      = New(http.StatusOK, 0, "success")
	ErrUnkown       = New(http.StatusOK, -1, "unknown error")
	ErrServicePanic = New(http.StatusServiceUnavailable, -2, "service panic")

	ErrCtxCanceled         = New(http.StatusOK, -11, "context canceled")
	ErrCtxDeadlineExceeded = New(http.StatusOK, -12, "context deadline exceeded")

	ErrReflectError        = New(http.StatusOK, -21, "reflect error")
	ErrReflectTargetNotPtr = New(http.StatusOK, -22, "reflect target is not pointer")
	ErrReflectTargetNilPtr = New(http.StatusOK, -23, "reflect target is nil pointer")
	ErrReflectTypeNotEqual = New(http.StatusOK, -24, "reflect types of target and source are not equal")

	ErrNonAuthoritativeInfo = New(http.StatusNonAuthoritativeInfo, -101, "http non authoritative info")
	ErrUnauthorized         = New(http.StatusUnauthorized, -102, "http unauthorized")
	ErrInternalServerError  = New(http.StatusInternalServerError, -103, "http internal server error")
	ErrTooManyRequests      = New(http.StatusTooManyRequests, -104, "too many requests")

	ErrInvalidRequestType = New(http.StatusOK, -201, "input request type must be pointer to struct")
	ErrParseRequest       = New(http.StatusOK, -202, "parse request error")
	ErrRequestValidation  = New(http.StatusOK, -203, "request validation error")
	ErrDuplicateRequest   = New(http.StatusOK, -204, "duplicate request")

	ErrHttpCall          = New(http.StatusOK, -301, "http call request error")
	ErrHttpReadBody      = New(http.StatusOK, -302, "http read body error")
	ErrHttpNewRequest    = New(http.StatusOK, -303, "http new request error")
	ErrHttpInvalidStatus = New(http.StatusOK, -304, "http invalid status")
	ErrHttpRenderBody    = New(http.StatusOK, -305, "http render body fail")

	ErrRpcCall      = New(http.StatusOK, -311, "rpc call error")
	ErrRpcTimeout   = New(http.StatusOK, -312, "rpc timeout")
	ErrRpcBokenPipe = New(http.StatusOK, -313, "rpc boken pip")

	ErrDataNotFound  = New(http.StatusOK, -801, "data not found")
	ErrDuplicateData = New(http.StatusOK, -802, "duplicate data")

	ErrDbQuery    = New(http.StatusOK, -811, "database query error")
	ErrDbExcute   = New(http.StatusOK, -812, "database execute error")
	ErrDbTxBegin  = New(http.StatusOK, -813, "database tx begin error")
	ErrDbTxCommit = New(http.StatusOK, -814, "database tx commit error")

	ErrCacheRead      = New(http.StatusOK, -821, "cache read error")
	ErrCacheWrite     = New(http.StatusOK, -822, "cache write error")
	ErrCacheTimeout   = New(http.StatusOK, -823, "cache timeout")
	ErrCacheMarshal   = New(http.StatusOK, -824, "cache codec marshal error")
	ErrCacheUnmarshal = New(http.StatusOK, -825, "cache codec unmarshal error")

	ErrCacheMutexLocked    = New(http.StatusOK, -831, "cache mutex had been locked")
	ErrCacheMutexNotExists = New(http.StatusOK, -832, "cache mutex is not exists")
	ErrCacheMutexNotMatch  = New(http.StatusOK, -833, "cache mutex is not match")
)

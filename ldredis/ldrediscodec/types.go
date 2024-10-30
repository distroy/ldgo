/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"github.com/distroy/ldgo/v3/ldredis"
)

type (
	client   = ldredis.Redis
	Reporter = ldredis.Reporter
)

type ZMember[T any] struct {
	Score  float64
	Member T
}

type ZMemberWithKey[T any] struct {
	ZMember[T]
	Key string
}

type ZAddArgs[T any] struct {
	NX      bool
	XX      bool
	LT      bool
	GT      bool
	Ch      bool
	Members []ZMember[T]
}

type Pair[T any] struct {
	First  string // key or field
	Second T      // value
}

type Cmdable[T comparable] interface {
	StringCmdable[T]
	HashCmdable[T]
	ListCmdable[T]
	SetCmdable[T]
	SortedSetCmdable[T]

	// XAdd(a *XAddArgs) *StringCmd
}

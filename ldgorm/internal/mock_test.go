/*
 * Copyright (C) distroy
 */

package internal

import (
	"fmt"
	"sync/atomic"

	"github.com/distroy/ldgo/v3/ldhook"
	"gorm.io/gorm"
)

var (
	testErrTransaction = fmt.Errorf("transaction error")
)

type transactionMockData struct {
	Begin    int32
	Commit   int32
	Rollback int32
	// RollbackUnlessCommitted int32
}

func mockTransaction(patches ldhook.Patches) *transactionMockData {
	data := &transactionMockData{}

	db := &gorm.DB{}
	patches.Applys([]ldhook.Hook{
		ldhook.FuncHook{
			Target: (*gorm.DB).Begin,
			Double: ldhook.ResultCell{
				Outputs: ldhook.Values{db},
				Inputs: ldhook.Values{
					ldhook.BindInput(0, func(db *gorm.DB) {
						atomic.AddInt32(&data.Begin, 1)
					}),
				},
			},
		},
		ldhook.FuncHook{
			Target: (*gorm.DB).Commit,
			Double: ldhook.ResultCell{
				Outputs: ldhook.Values{db},
				Inputs: ldhook.Values{
					ldhook.BindInput(0, func(db *gorm.DB) {
						atomic.AddInt32(&data.Commit, 1)
					}),
				},
			},
		},
		ldhook.FuncHook{
			Target: (*gorm.DB).Rollback,
			Double: ldhook.ResultCell{
				Outputs: ldhook.Values{db},
				Inputs: ldhook.Values{
					ldhook.BindInput(0, func(db *gorm.DB) {
						atomic.AddInt32(&data.Rollback, 1)
					}),
				},
			},
		},
	})

	return data
}

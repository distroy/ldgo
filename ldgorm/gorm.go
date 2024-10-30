/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"fmt"

	"github.com/distroy/ldgo/v3/ldgorm/internal"
	"gorm.io/gorm"
)

type (
	OriginGormDb = gorm.DB
	GormDb       = internal.GormDb
)

func New(db *gorm.DB) *GormDb {
	return internal.New(db)
}

func NewDb(db *gorm.DB) *GormDb {
	return New(db)
}

func NewTestDb() (*GormDb, error) {
	return internal.NewTestGormDb()
}

func MustNewTestDb() *GormDb {
	db, err := NewTestDb()
	if err != nil {
		panic(fmt.Sprintf("new test gorm fail. err:%v", err))
	}

	return db
}

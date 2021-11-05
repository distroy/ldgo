/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"fmt"

	"github.com/distroy/ldgo/ldmath"
)

type Option interface {
	buildGorm(db *GormDb) *GormDb
	String() string
}

func ApplyOptions(db *GormDb, opts ...Option) *GormDb {
	for _, opt := range opts {
		db = opt.buildGorm(db)
	}
	return db
}

type pagingOption struct {
	Page     int64 // the first page is 1
	Pagesize int64
}

func (p pagingOption) String() string {
	return fmt.Sprintf(`{"page":%d,"pagesize":%d}`, p.Page, p.Pagesize)
}

func (p pagingOption) buildGorm(db *GormDb) *GormDb {
	if p.Pagesize > 0 {
		p.Page = ldmath.MaxInt64(1, p.Page)
		offset := (p.Page - 1) * p.Pagesize
		db = db.Offset(offset).Limit(p.Pagesize)
	}

	return db
}

// Paging return the paging option
// the first page is 1
// if pagesize <= 0, it will query all rows
func Paging(page int64, pagesize int64) Option {
	return pagingOption{
		Page:     page,
		Pagesize: pagesize,
	}
}

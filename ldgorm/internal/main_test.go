/*
 * Copyright (C) distroy
 */

package internal

import "github.com/distroy/ldgo/v3/ldlog"

type testTable struct {
	ProjectId int64 `gorm:"column:project_id;primary_key"`
	ChannelId int64 `gorm:"column:channel_id"`
	Type      int64 `gorm:"column:type"`
	VersionId int64 `gorm:"column:version_id"`
}

func (_ *testTable) TableName() string { return "test_table" }

func testNewGorm() *GormDb {
	db, _ := NewTestGormDb()
	db = db.WithLogger(ldlog.Discard().Wrapper())
	db.CreateTable(&testTable{})
	return db
}

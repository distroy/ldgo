/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"strings"
	"testing"

	"github.com/distroy/ldgo/ldlogger"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

type testFilter struct {
	VersionId FieldWherer `gormwhere:"column:version_id;"`
	ChannelId FieldWherer `gormwhere:"column:channel_id;order:2;notempty"`
	ProjectId FieldWherer `gormwhere:"column:project_id;order:1"`
	Type      FieldWherer `gormwhere:"column:type;"`
}

type testTable struct {
	ProjectId int64 `gorm:"column:project_id"`
	ChannelId int64 `gorm:"column:channel_id"`
	Type      int64 `gorm:"column:type"`
	VersionId int64 `gorm:"column:version_id"`
}

func (_ *testTable) TableName() string { return "test_table" }

func testGetGorm() *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	// convey.So(err, convey.ShouldBeNil)
	db.LogMode(false)
	db.CreateTable(&testTable{})
	db.SetLogger(ldlogger.Console().WithOptions(zap.IncreaseLevel(zap.ErrorLevel)).Wrap())
	return db
}

func testGetWhereFromSql(scope *gorm.Scope) string {
	const key = " WHERE "
	sql := scope.SQL
	idx := strings.Index(sql, key)
	if idx < 0 {
		return ""
	}
	return sql[idx+len(key):]
}

func TestWhereOption(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		gormDb := testGetGorm()
		defer gormDb.Close()

		var res whereResult
		gormDb.Callback().Query().After("gorm:query").Register("ldgorm:after_query", func(scope *gorm.Scope) {
			res.Query = testGetWhereFromSql(scope)
			res.Args = scope.SQLVars
		})

		var rows []*testTable

		convey.Convey("(project_id = 100 && channel_id > 100) || (project_id = 123 && channel_id < 234)", func() {
			where := Where(&testFilter{
				ProjectId: Equal(10),
				ChannelId: Gt(100),
			}).Or(&testFilter{
				ProjectId: Equal(123),
				ChannelId: Lt(234),
			})

			where.buildGorm(gormDb).Find(&rows)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: "((`project_id` = ? AND `channel_id` > ?) OR (`project_id` = ? AND `channel_id` < ?))",
				Args:  []interface{}{10, 100, 123, 234},
			})
		})

		convey.Convey("(project_id = 10 && channel_id >= 0) && ((channel_id < 100 && version_id > 220) || channel_id > 200 && version_id < 110)", func() {
			where1 := Where(&testFilter{
				ProjectId: Equal(10),
				ChannelId: Between(0, nil),
			})
			where2 := Where(&testFilter{
				ChannelId: Lt(100),
				VersionId: Gt(220),
			}).Or(&testFilter{
				ChannelId: Gt(200),
				VersionId: Lt(110),
			})
			where := where1.And(where2)

			where.buildGorm(gormDb).Find(&rows)
			convey.So(res, convey.ShouldResemble, whereResult{
				Query: "((`project_id` = ? AND `channel_id` >= ?) AND ((`channel_id` < ? AND `version_id` > ?) OR (`channel_id` > ? AND `version_id` < ?)))",
				Args:  []interface{}{10, 0, 100, 220, 200, 110},
			})
		})
	})
}
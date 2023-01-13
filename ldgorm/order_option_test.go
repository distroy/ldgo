/*
 * Copyright (C) distroy
 */

package ldgorm

import (
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

type testOrderStruct struct {
	ProjectId FieldOrderer `gormorder:"column:project_id"`
	ChannelId FieldOrderer `gormorder:"column:channel_id"`
	VersionId FieldOrderer `gormorder:"column:version_id"`
	Type      FieldOrderer `gormorder:"column:type"`
}

func testGetOrderFromSql(scope *gorm.Scope) string {
	const key = " ORDER BY "
	sql := scope.SQL
	idx := strings.Index(sql, key)
	if idx < 0 {
		return ""
	}
	return sql[idx+len(key):]
}

func TestOrder(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		gormDb := testGetGorm()
		defer gormDb.Close()

		var res string
		var args []interface{}
		gormDb.Callback().Query().After("gorm:query").Register("ldgorm:after_query", func(scope *gorm.Scope) {
			res = testGetOrderFromSql(scope)
			args = scope.SQLVars
		})

		var rows []*testTable

		convey.Convey("channel_id, version_id DESC", func() {
			order := Order(&testOrderStruct{
				ChannelId: FieldOrder(1).Asc(),
				VersionId: FieldOrder(2).Desc(),
			})

			ApplyOptions(gormDb, order).Find(&rows)
			convey.So(res, convey.ShouldResemble, "`channel_id`,`version_id` DESC")
		})

		convey.Convey("channel_id DESC, type", func() {
			order := Order(&testOrderStruct{
				ChannelId: FieldOrder(1).Desc(),
				Type:      FieldOrder(2),
			})

			ApplyOptions(gormDb, order).Find(&rows)
			convey.So(res, convey.ShouldResemble, "`channel_id` DESC,`type`")
		})

		convey.Convey("FIELD type", func() {
			order := Order(&testOrderStruct{
				Type: FieldOrder(1).Field([]int{2, 4, 3}),
			})

			ApplyOptions(gormDb, order).Find(&rows)
			convey.So(res, convey.ShouldResemble, "FIELD(`type`, ?, ?, ?)")
			convey.So(args, convey.ShouldResemble, []interface{}{2, 4, 3})
		})

		convey.Convey("FIELD type DESC", func() {
			order := Order(&testOrderStruct{
				Type: FieldOrder(1).Field([]int{2, 4, 3}).Desc(),
			})

			ApplyOptions(gormDb, order).Find(&rows)
			convey.So(res, convey.ShouldResemble, "FIELD(`type`, ?, ?, ?) DESC")
			convey.So(args, convey.ShouldResemble, []interface{}{2, 4, 3})
		})
	})
}

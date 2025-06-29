/*
 * Copyright (C) distroy
 */

package internal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/distroy/ldgo/v3/ldhook"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/distroy/ldgo/v3/ldrand"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func TestGormDb_Transaction(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()

		data := mockTransaction(patches)

		db := testNewGorm()
		defer db.Close()

		convey.Convey("transaction 1 rollback", func() {
			err := db.Transaction(func(tx *GormDb) error {
				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("transaction 2 commit", func() {
					err := tx.Transaction(func(tx *GormDb) error {
						convey.So(data, convey.ShouldResemble, &transactionMockData{
							Begin: 1,
						})
						return nil
					})

					convey.So(err, convey.ShouldBeNil)
				})

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				return testErrTransaction
			})

			convey.So(err, convey.ShouldEqual, testErrTransaction)
			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:    1,
				Rollback: 1,
			})
		})

		convey.Convey("transaction 1 commit", func() {
			err := db.Transaction(func(tx *GormDb) error {
				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				return nil
			})

			convey.So(err, convey.ShouldBeNil)
			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:  1,
				Commit: 1,
			})
		})

		convey.Convey("transaction panic", func() {
			convey.So(func() {
				db.Transaction(func(tx *GormDb) error {
					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
					panic("transaction panic")
				})
			}, convey.ShouldPanic)
			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:    1,
				Rollback: 1,
			})
		})
	})
}

func TestGormDb_Begin_Commit(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()

		data := mockTransaction(patches)

		db := &GormDb{gormDb: &gorm.DB{}}

		convey.Convey("transaction 1 commit", func() {
			tx1Begin := db.Begin()

			convey.So(tx1Begin.Error, convey.ShouldBeNil)
			convey.So(tx1Begin.inTx, convey.ShouldEqual, true)
			convey.So(tx1Begin.txLvl, convey.ShouldEqual, 1)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin: 1,
			})

			convey.Convey("transaction 2 commit", func() {
				tx2Begin := tx1Begin.Begin()

				convey.So(tx2Begin.Error, convey.ShouldBeNil)
				convey.So(tx2Begin.inTx, convey.ShouldEqual, true)
				convey.So(tx2Begin.txLvl, convey.ShouldEqual, 2)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				tx2Commit := tx2Begin.Commit()

				convey.So(tx2Commit.Error, convey.ShouldBeNil)
				convey.So(tx2Commit.inTx, convey.ShouldEqual, false)
				convey.So(tx2Commit.txLvl, convey.ShouldEqual, 1)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("double commit", func() {
					convey.So(func() { tx2Commit.Commit() }, convey.ShouldPanic)
					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})

				convey.Convey("double rollback unless committed", func() {
					tx2RollbackUnlessCommitted := tx2Commit.RollbackUnlessCommitted()

					convey.So(tx2RollbackUnlessCommitted.Error, convey.ShouldBeNil)
					convey.So(tx2RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
					convey.So(tx2RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 1)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})
			})

			tx1Commit := tx1Begin.Commit()

			convey.So(tx1Commit.Error, convey.ShouldBeNil)
			convey.So(tx1Commit.inTx, convey.ShouldEqual, false)
			convey.So(tx1Commit.txLvl, convey.ShouldEqual, 0)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:  1,
				Commit: 1,
			})

			convey.Convey("double commit", func() {
				convey.So(func() { tx1Commit.Commit() }, convey.ShouldPanic)
				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:  1,
					Commit: 1,
				})
			})

			convey.Convey("double rollback unless committed", func() {
				tx1RollbackUnlessCommitted := tx1Commit.RollbackUnlessCommitted()

				convey.So(tx1RollbackUnlessCommitted.Error, convey.ShouldBeNil)
				convey.So(tx1RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
				convey.So(tx1RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 0)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:  1,
					Commit: 1,
				})
			})
		})
	})
}

func TestGormDb_Begin_Rollback(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()

		data := mockTransaction(patches)

		db := &GormDb{gormDb: &gorm.DB{}}

		convey.Convey("transaction 1 rollback", func() {
			tx1Begin := db.Begin()

			convey.So(tx1Begin.Error, convey.ShouldBeNil)
			convey.So(tx1Begin.inTx, convey.ShouldEqual, true)
			convey.So(tx1Begin.txLvl, convey.ShouldEqual, 1)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin: 1,
			})

			convey.Convey("transaction 2 rollback", func() {
				tx2Begin := tx1Begin.Begin()

				convey.So(tx2Begin.Error, convey.ShouldBeNil)
				convey.So(tx2Begin.inTx, convey.ShouldEqual, true)
				convey.So(tx2Begin.txLvl, convey.ShouldEqual, 2)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("transaction 3 commit", func() {
					tx3Begin := tx2Begin.Begin()

					convey.So(tx3Begin.Error, convey.ShouldBeNil)
					convey.So(tx3Begin.inTx, convey.ShouldEqual, true)
					convey.So(tx3Begin.txLvl, convey.ShouldEqual, 3)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})

					tx3Commit := tx3Begin.Commit()

					convey.So(tx3Commit.Error, convey.ShouldBeNil)
					convey.So(tx3Commit.inTx, convey.ShouldEqual, false)
					convey.So(tx3Commit.txLvl, convey.ShouldEqual, 2)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})

				tx2Rollback := tx2Begin.Rollback()

				convey.So(tx2Rollback.Error, convey.ShouldBeNil)
				convey.So(tx2Rollback.inTx, convey.ShouldEqual, false)
				convey.So(tx2Rollback.txLvl, convey.ShouldEqual, 1)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("double rollback", func() {
					convey.So(func() { tx2Rollback.Rollback() }, convey.ShouldPanic)
					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})

				convey.Convey("double rollback unless committed", func() {
					tx2RollbackUnlessCommitted := tx2Rollback.RollbackUnlessCommitted()

					convey.So(tx2RollbackUnlessCommitted.Error, convey.ShouldBeNil)
					convey.So(tx2RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
					convey.So(tx2RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 1)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})
			})

			tx1Rollback := tx1Begin.Rollback()

			convey.So(tx1Rollback.Error, convey.ShouldBeNil)
			convey.So(tx1Rollback.inTx, convey.ShouldEqual, false)
			convey.So(tx1Rollback.txLvl, convey.ShouldEqual, 0)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:    1,
				Rollback: 1,
			})

			convey.Convey("double rollback", func() {
				convey.So(func() { tx1Rollback.Rollback() }, convey.ShouldPanic)
				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:    1,
					Rollback: 1,
				})
			})

			convey.Convey("double rollback unless committed", func() {
				tx1RollbackUnlessCommitted := tx1Rollback.RollbackUnlessCommitted()

				convey.So(tx1RollbackUnlessCommitted.Error, convey.ShouldBeNil)
				convey.So(tx1RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
				convey.So(tx1RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 0)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:    1,
					Rollback: 1,
				})
			})
		})
	})
}

func TestGormDb_Begin_RollbackUnlessCommitted(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()

		data := mockTransaction(patches)

		db := &GormDb{gormDb: &gorm.DB{}}

		convey.Convey("transaction 1 rollback", func() {
			tx1Begin := db.Begin()

			convey.So(tx1Begin.Error, convey.ShouldBeNil)
			convey.So(tx1Begin.inTx, convey.ShouldEqual, true)
			convey.So(tx1Begin.txLvl, convey.ShouldEqual, 1)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin: 1,
			})

			convey.Convey("transaction 2 rollback", func() {
				tx2Begin := tx1Begin.Begin()

				convey.So(tx2Begin.Error, convey.ShouldBeNil)
				convey.So(tx2Begin.inTx, convey.ShouldEqual, true)
				convey.So(tx2Begin.txLvl, convey.ShouldEqual, 2)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("transaction 3 commit", func() {
					tx3Begin := tx2Begin.Begin()

					convey.So(tx3Begin.Error, convey.ShouldBeNil)
					convey.So(tx3Begin.inTx, convey.ShouldEqual, true)
					convey.So(tx3Begin.txLvl, convey.ShouldEqual, 3)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})

					tx3Commit := tx3Begin.Commit()

					convey.So(tx3Commit.Error, convey.ShouldBeNil)
					convey.So(tx3Commit.inTx, convey.ShouldEqual, false)
					convey.So(tx3Commit.txLvl, convey.ShouldEqual, 2)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})

				tx2Rollback := tx2Begin.RollbackUnlessCommitted()

				convey.So(tx2Rollback.Error, convey.ShouldBeNil)
				convey.So(tx2Rollback.inTx, convey.ShouldEqual, false)
				convey.So(tx2Rollback.txLvl, convey.ShouldEqual, 1)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin: 1,
				})

				convey.Convey("double rollback", func() {
					convey.So(func() { tx2Rollback.Rollback() }, convey.ShouldPanic)
					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})

				convey.Convey("double rollback unless committed", func() {
					tx2RollbackUnlessCommitted := tx2Rollback.RollbackUnlessCommitted()

					convey.So(tx2RollbackUnlessCommitted.Error, convey.ShouldBeNil)
					convey.So(tx2RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
					convey.So(tx2RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 1)

					convey.So(data, convey.ShouldResemble, &transactionMockData{
						Begin: 1,
					})
				})
			})

			tx1Rollback := tx1Begin.RollbackUnlessCommitted()

			convey.So(tx1Rollback.Error, convey.ShouldBeNil)
			convey.So(tx1Rollback.inTx, convey.ShouldEqual, false)
			convey.So(tx1Rollback.txLvl, convey.ShouldEqual, 0)

			convey.So(data, convey.ShouldResemble, &transactionMockData{
				Begin:    1,
				Rollback: 1,
			})

			convey.Convey("double rollback", func() {
				convey.So(func() { tx1Rollback.Rollback() }, convey.ShouldPanic)
				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:    1,
					Rollback: 1,
				})
			})

			convey.Convey("double rollback unless committed", func() {
				tx1RollbackUnlessCommitted := tx1Rollback.RollbackUnlessCommitted()

				convey.So(tx1RollbackUnlessCommitted.Error, convey.ShouldBeNil)
				convey.So(tx1RollbackUnlessCommitted.inTx, convey.ShouldEqual, false)
				convey.So(tx1RollbackUnlessCommitted.txLvl, convey.ShouldEqual, 0)

				convey.So(data, convey.ShouldResemble, &transactionMockData{
					Begin:    1,
					Rollback: 1,
				})
			})
		})
	})
}

func TestGormDb_WithQueryHint(t *testing.T) {
	tests := []struct {
		name string
		args []func(db *GormDb) *GormDb
		want bool
	}{
		{
			name: "select",
			want: true,
			args: []func(db *GormDb) *GormDb{
				func(db *GormDb) *GormDb { return db.Find(&[]*testTable{}) },
			},
		},
		{
			name: "select limit 1",
			want: true,
			args: []func(db *GormDb) *GormDb{
				func(db *GormDb) *GormDb { return db.First(&testTable{}) },
			},
		},
		{
			name: "insert",
			want: false,
			args: []func(db *GormDb) *GormDb{
				func(db *GormDb) *GormDb { return db.Create(&testTable{ProjectId: 1}) },
			},
		},
		{
			name: "save",
			want: false,
			args: []func(db *GormDb) *GormDb{
				func(db *GormDb) *GormDb { return db.Save(&testTable{ProjectId: 1}) },
			},
		},
		{
			name: "update",
			want: false,
			args: []func(db *GormDb) *GormDb{
				func(db *GormDb) *GormDb { return db.Updates(&testTable{ProjectId: 1}) },
			},
		},
	}

	hint := ldrand.String(8)
	prefix := fmt.Sprintf("/* %s */", hint)

	db := testNewGorm()
	defer db.Close()
	db = db.WithQueryHint(hint)

	db.Save(&testTable{
		ProjectId: 1,
		ChannelId: 2,
		Type:      3,
		VersionId: 4,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *GormDb) *GormDb {
				db := tx
				for _, fn := range tt.args {
					db = fn(db)
				}
				return db
			})

			if got := strings.HasPrefix(sql, prefix); got != tt.want {
				t.Errorf("[FAIL] GormDb.WithQueryHint() = %v, want %v, sql = %s", got, tt.want, sql)
			} else {
				t.Logf("[PASS] GormDb.WithQueryHint() = %v, want %v, sql = %s", got, tt.want, sql)
			}
		})
	}
}

func TestGormDb_WithLogger(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("without slaver", func(c convey.C) {
			db := testNewGorm()
			defer db.Close()

			buf := &strings.Builder{}
			db = db.WithLogger(ldlog.New(buf).Wrapper())

			db = db.UseSlaver()

			db.Save(&testTable{
				ProjectId: 1001,
				ChannelId: 1002,
				Type:      1003,
				VersionId: 1004,
			})

			t.Logf("\n%s", buf.String())
			c.So(buf.String(), convey.ShouldNotEqual, ``)
		})

		c.Convey("with slaver", func(c convey.C) {
			c.Convey("rand", func(c convey.C) {
				db := testNewGorm()
				defer db.Close()

				db = db.AddSlaver(testNewGorm().Get())
				db = db.AddSlaver(testNewGorm().Get())

				db = db.UseSlaver()

				buf := &strings.Builder{}
				db = db.WithLogger(ldlog.New(buf).Wrapper())

				db.Save(&testTable{
					ProjectId: 1001,
					ChannelId: 1002,
					Type:      1003,
					VersionId: 1004,
				})

				t.Logf("\n%s", buf.String())
				c.So(buf.String(), convey.ShouldNotEqual, ``)
			})

			c.Convey("index", func(c convey.C) {
				db := testNewGorm()
				defer db.Close()

				db = db.AddSlaver(testNewGorm().Get())
				db = db.AddSlaver(testNewGorm().Get())

				db = db.UseSlaver(1)

				buf := &strings.Builder{}
				db = db.WithLogger(ldlog.New(buf).Wrapper())

				db.Save(&testTable{
					ProjectId: 1001,
					ChannelId: 1002,
					Type:      1003,
					VersionId: 1004,
				})

				t.Logf("\n%s", buf.String())
				c.So(buf.String(), convey.ShouldNotEqual, ``)
			})
		})
	})
}

/*
 * Copyright (C) distroy
 */

package gorm

import (
	"github.com/distroy/ldgo/ldlog"
	"github.com/jinzhu/gorm"
)

type gormDb = gorm.DB

type GormDb struct {
	*gormDb

	txLvl int
	inTx  bool
}

func (w *GormDb) panicTxLevelLessZero() {
	panic("tx level must not be less than zero")
}

func (w *GormDb) panicTxCommittedOrRollbacked() {
	panic("tx can not be committed or rollbacked again")
}

func (w *GormDb) clone() *GormDb {
	c := *w
	return &c
}

// New clone a new db connection without search conditions
func (w *GormDb) New() *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.New()
	return w
}

func (w *GormDb) Get() *gorm.DB {
	return w.gormDb
}

func (w *GormDb) Set(db *gorm.DB) *GormDb {
	w = w.clone()
	w.gormDb = db
	return w
}

func (w *GormDb) WithLogger(l ldlog.LoggerInterface) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.LogMode(true)
	w.gormDb.SetLogger(ldlog.GetWrapper(l))
	return w
}

func (w *GormDb) Model(value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Model(value)
	return w
}

func (w *GormDb) Transaction(fc func(tx *GormDb) error) (err error) {
	if w.txLvl > 0 {
		return fc(w)
	}

	paniced := true
	tx := w.Begin()
	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if paniced || err != nil {
			tx.Rollback()
		}
	}()

	err = fc(tx)

	if err == nil {
		err = tx.Commit().Error
	}

	paniced = false
	return
}

func (w *GormDb) Begin() *GormDb {
	if w.txLvl < 0 {
		w.panicTxLevelLessZero()
	}

	w = w.clone()

	if w.txLvl == 0 {
		w.gormDb = w.gormDb.Begin()
	}

	w.inTx = true
	w.txLvl++
	return w
}

func (w *GormDb) Commit() *GormDb {
	if !w.inTx {
		w.panicTxCommittedOrRollbacked()
	}

	w = w.clone()

	w.txLvl--
	if w.txLvl < 0 {
		w.panicTxLevelLessZero()
	}

	if w.txLvl == 0 {
		w.gormDb = w.gormDb.Commit()
	}

	w.inTx = false
	return w
}

func (w *GormDb) Rollback() *GormDb {
	if !w.inTx {
		w.panicTxCommittedOrRollbacked()
	}

	w = w.clone()

	w.txLvl--
	if w.txLvl < 0 {
		w.panicTxLevelLessZero()
	}

	if w.txLvl == 0 {
		w.gormDb = w.gormDb.Rollback()
	}

	w.inTx = false
	return w
}

func (w *GormDb) RollbackUnlessCommitted() *GormDb {
	if !w.inTx {
		return w
	}

	w = w.clone()

	w.txLvl--
	if w.txLvl < 0 {
		w.panicTxLevelLessZero()
	}

	if w.txLvl == 0 {
		w.gormDb = w.gormDb.RollbackUnlessCommitted()
	}

	w.inTx = false
	return w
}

func (w *GormDb) Select(query interface{}, args ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Select(query, args...)
	return w
}

func (w *GormDb) Group(query string) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Group(query)
	return w
}

func (w *GormDb) Having(query interface{}, args ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Having(query, args...)
	return w
}

func (w *GormDb) Joins(query string, args ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Joins(query, args...)
	return w
}

func (w *GormDb) HasTable(value interface{}) bool {
	return w.gormDb.HasTable(value)
}

func (w *GormDb) CreateTable(models ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.CreateTable(models...)
	return w
}

func (w *GormDb) DropTable(models ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.DropTable(models...)
	return w
}

func (w *GormDb) DropTableIfExists(models ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.DropTableIfExists(models...)
	return w
}

func (w *GormDb) Where(query interface{}, args ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Where(query, args...)
	return w
}

// Order specify order when retrieve records from database, set reorder to `true` to overwrite defined conditions
//     db.Order("name DESC")
//     db.Order("name DESC", true) // reorder
//     db.Order(gorm.Expr("name = ? DESC", "first")) // sql expression
func (w *GormDb) Order(value interface{}, reorder ...bool) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Order(value, reorder...)
	return w
}

func (w *GormDb) Limit(limit interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Limit(limit)
	return w
}

func (w *GormDb) Offset(offset interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Offset(offset)
	return w
}

func (w *GormDb) Save(value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Save(value)
	return w
}

func (w *GormDb) Create(value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Create(value)
	return w
}

func (w *GormDb) Update(attrs ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Update(attrs...)
	return w
}

func (w *GormDb) Updates(values interface{}, ignoreProtectedAttrs ...bool) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Updates(values, ignoreProtectedAttrs...)
	return w
}

func (w *GormDb) FirstOrCreate(out interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.FirstOrCreate(out, where...)
	return w
}

func (w *GormDb) Delete(value interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Delete(value, where...)
	return w
}

func (w *GormDb) First(out interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.First(out, where...)
	return w
}

func (w *GormDb) Find(out interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Find(out, where...)
	return w
}

func (w *GormDb) Take(out interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Take(out, where...)
	return w
}

func (w *GormDb) Last(out interface{}, where ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Last(out, where...)
	return w
}

func (w *GormDb) Count(out interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Count(out)
	return w
}

func (w *GormDb) Exec(sql string, values ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Exec(sql, values...)
	return w
}

func (w *GormDb) Raw(sql string, values ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Raw(sql, values...)
	return w
}

func (w *GormDb) Scan(out interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Scan(out)
	return w
}

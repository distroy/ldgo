/*
 * Copyright (C) distroy
 */

package internal

/*
 * Filename: gorm.gen.go
 * The file name suffix ".gen.go" is used to prevent gorm from printing logs and to redirect the path to this file.
 */

import (
	"hash/crc32"
	"strings"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldrand"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/hints"
)

type gormDb = gorm.DB

type Logger interface {
	Printf(fmt string, v ...interface{})
}

var (
	queryHintReplacer = strings.NewReplacer("/*", " ", "*/", " ")
)

func New(db *gorm.DB) *GormDb {
	g := &GormDb{}
	g = g.Set(db)
	return g
}

func NewTestGormDb() (*GormDb, error) {
	// db, err := gorm.Open("sqlite3", ":memory:")
	const InMemoryDSN = "file:testdatabase?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(InMemoryDSN))
	if err != nil {
		return nil, err
	}

	return New(db), nil
}

type GormDb struct {
	*gormDb // it is currently used db, it is the master by default

	master  *gorm.DB
	slavers []*gorm.DB

	txLvl int
	inTx  bool

	// these should set after user new gorm db
	log logger.Interface
	ctx ldctx.Context
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
	w.gormDb = w.gormDb.Session(&gorm.Session{
		Logger: w.gormDb.Logger,
	})
	return w
}

func (w *GormDb) Get() *gorm.DB {
	return w.gormDb
}

func (w *GormDb) Close() error {
	var err error

	db := w.master
	if e := w.close(db); e != nil {
		err = e
	}

	for _, db := range w.slavers {
		if e := w.close(db); err == nil && e != nil {
			err = e
		}
	}

	return err
}

func (w *GormDb) close(db *gorm.DB) error {
	x, err := db.DB()
	if err != nil {
		return err
	}
	return x.Close()
}

// Set will set the db currently used.
// If master is not set, will set the master db also.
func (w *GormDb) Set(db *gorm.DB) *GormDb {
	w = w.clone()

	w.gormDb = db
	if w.master == nil {
		w.master = db
	}

	return w
}

func (w *GormDb) SetMaster(db *gorm.DB) *GormDb {
	w = w.clone()

	w.master = db
	if w.gormDb == nil {
		w.gormDb = db
	}

	return w
}

func (w *GormDb) AddSlaver(dbs ...*gorm.DB) *GormDb {
	if len(dbs) == 0 {
		return w
	}

	w = w.clone()

	slavers := make([]*gorm.DB, 0, len(w.slavers)+len(dbs))
	slavers = append(slavers, w.slavers...)
	slavers = append(slavers, dbs...)
	w.slavers = slavers

	if w.gormDb == nil {
		w.gormDb = dbs[0]
	}
	if w.master == nil {
		w.master = dbs[0]
	}

	return w
}

// UseMaster must be called before all Create/Update/Query/Delete methods
func (w *GormDb) UseMaster() *GormDb {
	if w.gormDb == w.master {
		return w
	}

	w = w.clone()
	w.gormDb = w.master
	w.initAfterUseNewGormDb()
	return w
}

// UseSlaver must be called before all Query methods
//
// key must be int/int{8-64}/uint/uint{8-64}/uintptr/string/[]byte.
// if key is not set, will use rand slaver
func (w *GormDb) UseSlaver(key ...interface{}) *GormDb {
	n := len(w.slavers)
	switch n {
	case 0:
		return w

	case 1:
		w = w.clone()
		w.gormDb = w.slavers[0]
		return w
	}

	hash := w.getHashByKey(key)
	idx := hash % uint(n)

	w = w.clone()
	w.gormDb = w.slavers[idx]
	w.initAfterUseNewGormDb()
	return w
}

func (w *GormDb) initAfterUseNewGormDb() {
	db := w.gormDb

	db = w.withContext(w.ctx, db)
	db = w.withLogger(db, w.log)

	w.gormDb = db
}

func (w *GormDb) getHashByKey(keys []interface{}) uint {
	if len(keys) == 0 {
		return ldrand.Uint()
	}

	switch v := keys[0].(type) {
	case int:
		return uint(v)
	case int8:
		return uint(v)
	case int16:
		return uint(v)
	case int32:
		return uint(v)
	case int64:
		return uint(v)

	case uint:
		return uint(v)
	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case uintptr:
		return uint(v)

	case string:
		return uint(crc32.ChecksumIEEE(ldconv.StrToBytesUnsafe(v)))

	case []byte:
		return uint(crc32.ChecksumIEEE(v))
	}

	return ldrand.Uint()
}

// WithContext can be called before or after UseMaster/UseSlaver
func (w *GormDb) WithContext(ctx ldctx.Context) *GormDb {
	w = w.clone()

	w.ctx = ctx
	w.log = nil

	w.gormDb = w.withContext(ctx, w.gormDb)
	return w
}

func (_ *GormDb) withContext(ctx ldctx.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}

	writer := ldctx.GetLogger(ctx).Wrapper()
	return db.Session(&gorm.Session{
		Logger:  logger.New(writer, getLoggerConfig()),
		Context: ctx,
	})
}

// WithLogger can be called before or after UseMaster/UseSlaver
func (w *GormDb) WithLogger(l Logger) *GormDb {
	w = w.clone()

	log := logger.New(l, getLoggerConfig())
	w.log = log

	w.gormDb = w.withLogger(w.gormDb, log)
	return w
}

func (_ *GormDb) withLogger(db *gorm.DB, log logger.Interface) *gorm.DB {
	if log == nil {
		return db
	}
	return db.Session(&gorm.Session{
		Logger: log,
	})
}

func (w *GormDb) Model(value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Model(value)
	return w
}

func (w *GormDb) Table(table string) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Table(table)
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
		err := w.gormDb.Error
		w.gormDb = w.gormDb.Rollback()
		w.gormDb.Error = err
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

func (w *GormDb) CreateTable(models ...interface{}) error {
	w = w.clone()
	return w.gormDb.AutoMigrate(models...)
}

func (w *GormDb) Clauses(conds ...clause.Expression) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Clauses(conds...)
	return w
}

func (w *GormDb) Where(query interface{}, args ...interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Where(query, args...)
	return w
}

// Order specify order when retrieve records from database, set reorder to `true` to overwrite defined conditions
//
//	db.Order("name DESC")
func (w *GormDb) Order(value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Order(value)
	return w
}

func (w *GormDb) Limit(limit int) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Limit(limit)
	return w
}

func (w *GormDb) Offset(offset int) *GormDb {
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

func (w *GormDb) Update(column string, value interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Update(column, value)
	return w
}

func (w *GormDb) Updates(values interface{}) *GormDb {
	w = w.clone()
	w.gormDb = w.gormDb.Updates(values)
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

func (w *GormDb) Count(out *int64) *GormDb {
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

// WithQueryHint add hint comment before sql
//
// WithQueryHint must be called after UseSlaver
func (w *GormDb) WithQueryHint(hint string) *GormDb {
	w = w.clone()
	exprs := []clause.Expression{
		hints.CommentBefore("SELECT", hint),
		// hints.CommentBefore("UPDATE", hint),
		// hints.CommentBefore("INSERT", hint),
	}
	w.gormDb = w.gormDb.Clauses(exprs...)
	return w
}

func (w *GormDb) ToSQL(fn func(tx *GormDb) *GormDb) string {
	return w.gormDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		db := New(tx)
		db = fn(db)
		return db.Get()
	})
}

func getLoggerConfig() logger.Config {
	return logger.Config{
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		LogLevel:                  logger.Info,
	}
}

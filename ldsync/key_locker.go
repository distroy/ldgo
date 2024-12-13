/*
 * Copyright (C) distroy
 */

package ldsync

import "sync"

var (
	_ Locker = (*lockerPoolLocker)(nil)
)

var lockerPoolDataPool = &Pool[*lockerPoolData]{
	New: func() *lockerPoolData { return &lockerPoolData{} },
}

func NewLockerPool() *LockerPool {
	return &LockerPool{}
}

type lockerPoolData struct {
	locker sync.Mutex
	count  int64
}

type LockerPool struct {
	locker sync.Mutex
	keys   map[interface{}]*lockerPoolData
}

func (kl *LockerPool) init() {
	if kl.keys == nil {
		kl.keys = make(map[interface{}]*lockerPoolData, 1)
	}
}

func (kl *LockerPool) Get(key interface{}) Locker {
	var l Locker = &lockerPoolLocker{
		pool: kl,
		key:  key,
	}
	// l = AutoLock(l)
	return l
}

func (kl *LockerPool) get(key interface{}) *lockerPoolData {
	kl.init()
	return kl.keys[key]
}

func (kl *LockerPool) new(key interface{}) *lockerPoolData {
	d := lockerPoolDataPool.Get()
	d.count = 0
	kl.keys[key] = d
	return d
}

type lockerPoolLocker struct {
	pool *LockerPool
	key  interface{}
	data *lockerPoolData
}

func (l *lockerPoolLocker) Lock() {
	kl := l.pool
	key := l.key

	if l.data != nil {
		return
	}

	kl.locker.Lock()
	d := kl.get(key)
	if d == nil {
		d = kl.new(key)
	}
	d.count++
	kl.locker.Unlock()

	d.locker.Lock()
	l.data = d
}

func (l *lockerPoolLocker) Unlock() {
	kl := l.pool
	key := l.key
	d := l.data

	if d == nil {
		return
	}

	kl.locker.Lock()
	// d := kl.getAndSub(key)
	d.count--
	if d.count <= 0 {
		delete(kl.keys, key)
		lockerPoolDataPool.Put(d)
	}
	kl.locker.Unlock()

	d.locker.Unlock()
	l.data = nil
}

func (l *lockerPoolLocker) TryLock() bool {
	kl := l.pool
	key := l.key

	if l.data != nil {
		return true
	}

	kl.locker.Lock()

	if d := kl.get(key); d != nil {
		kl.locker.Unlock()
		return false
	}

	d := kl.new(key)
	d.count++
	kl.locker.Unlock()

	l.data = d
	return d.locker.TryLock()
}

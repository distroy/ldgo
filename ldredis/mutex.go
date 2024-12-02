/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"encoding/hex"
	"time"
	"unsafe"

	"github.com/distroy/ldgo/v2/ldatomic"
	"github.com/distroy/ldgo/v2/ldctx"
	"github.com/distroy/ldgo/v2/lderr"
	"github.com/distroy/ldgo/v2/ldrand"
	"go.uber.org/zap"
)

type MutexEvent int

const (
	MutexEventDeleted MutexEvent = iota + 1
)

var (
	closedMutextEventsChan chan MutexEvent
)

const (
	mutexMinHeartbeatInterval = 1 * time.Second
	mutexMinHeartbeatTimeout  = 10 * time.Second
	mutexMinLockForceInterval = 1 * time.Millisecond
)

type mutexContext struct {
	ctx    ldctx.Context // control goroutine to exit
	cancel ldctx.CancelFunc

	key           string
	token         string
	lastHeartbeat time.Time
	lockTime      ldatomic.Int64 // if equal 0, has not locked
	events        chan MutexEvent
}

type atomicMutexContext struct {
	d ldatomic.Pointer
}

func (a *atomicMutexContext) toData(p unsafe.Pointer) *mutexContext {
	return (*mutexContext)(p)
}

func (a *atomicMutexContext) Load() *mutexContext   { return a.toData(a.d.Load()) }
func (a *atomicMutexContext) Store(d *mutexContext) { a.d.Store(unsafe.Pointer(d)) }

func (a *atomicMutexContext) Swap(new *mutexContext) (old *mutexContext) {
	return a.toData(a.d.Swap(unsafe.Pointer(new)))
}

func (a *atomicMutexContext) CompareAndSwap(old, new *mutexContext) (swapped bool) {
	return a.d.CompareAndSwap(unsafe.Pointer(old), unsafe.Pointer(new))
}

type Mutex struct {
	redis             *Redis
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration
	lockForce         bool
	lockForceInterval time.Duration
	lockForceTimeout  time.Duration
	unlockDelay       time.Duration

	// if equal nil, has not locked
	// but when mutex has been cloned. maybe the lockTime is equal 0, the ctx is not equal nil
	ctx atomicMutexContext
}

func NewMutex(redis *Redis) *Mutex {
	m := &Mutex{
		redis:             redis,
		heartbeatInterval: 5 * time.Second,
		heartbeatTimeout:  2 * time.Minute,
		lockForceInterval: 1 * time.Second,
	}
	return m
}

func (m *Mutex) clone() *Mutex {
	c := *m
	return &c
}

func (m *Mutex) WithRedis(redis *Redis) *Mutex {
	m = m.clone()
	m.redis = redis
	return m
}

// WithLockForce returns the redis mutex with lock force
// but if the context is cancelled, the lock force is not available
//
// WithLockForce should be called like these:
//
//	WithLockForce(false)
//	WithLockForce(true, {interval})
//	WithLockForce(true, {interval}, {timeout})
func (m *Mutex) WithLockForce(force bool, intervalAndTimeout ...time.Duration) *Mutex {
	m = m.clone()
	m.lockForce = force

	m.lockForceInterval = 0
	m.lockForceInterval = mutexMinLockForceInterval
	if len(intervalAndTimeout) > 0 {
		m.lockForceInterval = intervalAndTimeout[0]
	}
	if m.lockForceInterval < mutexMinLockForceInterval {
		m.lockForceInterval = mutexMinLockForceInterval
	}

	m.lockForceTimeout = 0
	if len(intervalAndTimeout) > 1 {
		m.lockForceTimeout = intervalAndTimeout[1]
	}

	return m
}

// WithUnlockDelay returns the redis mutex with unlock delay
func (m *Mutex) WithUnlockDelay(delay ...time.Duration) *Mutex {
	m = m.clone()
	m.unlockDelay = 0
	if len(delay) > 0 && delay[0] > 0 {
		m.unlockDelay = delay[0]
	}
	return m
}

func (m *Mutex) WithInterval(d time.Duration) *Mutex {
	m = m.clone()

	if d < mutexMinHeartbeatInterval {
		d = mutexMinHeartbeatInterval
	}
	m.heartbeatInterval = d

	if timeout := m.getMinTimeout(d); m.heartbeatTimeout < timeout {
		m.heartbeatTimeout = timeout
	}
	return m
}

func (m *Mutex) WithTimeout(d time.Duration) *Mutex {
	m = m.clone()
	m.heartbeatTimeout = m.getMinTimeout(d)
	return m
}

func (m *Mutex) getMinTimeout(d time.Duration) time.Duration {
	if d < mutexMinHeartbeatTimeout {
		d = mutexMinHeartbeatTimeout
	}

	if t := m.heartbeatInterval * 3; d < t {
		d = t
	}

	return d
}

func (m *Mutex) mustGetContext() *mutexContext {
	mc := m.ctx.Load()
	if mc == nil {
		mc = &mutexContext{
			events: closedMutextEventsChan,
		}
	}
	return mc
}

func (m *Mutex) getExpiration() time.Duration {
	return m.heartbeatInterval + m.heartbeatTimeout
}

func (m *Mutex) Key() string               { return m.mustGetContext().key }
func (m *Mutex) Events() <-chan MutexEvent { return m.mustGetContext().events }

func (m *Mutex) Lock(ctx ldctx.Context, key string) error {
	mc := m.ctx.Load()
	if mc != nil && mc.lockTime.Load() != 0 {
		ldctx.LogE(ctx, "redis mutex has been locked", zap.String("key", key), zap.String("old", mc.key),
			getCallerField(m.redis))
		return lderr.ErrCacheMutexLocked
	}
	if mc != nil {
		m.ctx.CompareAndSwap(mc, nil)
	}

	token := hex.EncodeToString(ldrand.Bytes(16))

	ctx, cancel := ldctx.WithCancel(ctx)
	ctx = ldctx.WithLogField(ctx, zap.String("key", key), zap.String("token", token))

	if err := m.internalLockOrLockForce(ctx, key, token); err != nil {
		return err
	}

	now := time.Now().UnixNano()
	mc = &mutexContext{
		ctx:      ctx,
		cancel:   cancel,
		key:      key,
		token:    token,
		lockTime: ldatomic.Int64(now),
		events:   make(chan MutexEvent, 1),
	}
	if ok := m.ctx.CompareAndSwap(nil, mc); !ok {
		// cli := m.redis.Client()
		// cli.Del(key)
		return lderr.ErrCacheMutexLocked
	}

	go m.goroutine(mc, now)

	ldctx.LogD(ctx, "redis mutex lock succ", getCallerField(m.redis))
	return nil
}

func (m *Mutex) internalLockOrLockForce(ctx ldctx.Context, key, token string) error {
	err := m.internalLock(ctx, key, token)
	if err == nil {
		return nil
	}

	if !m.lockForce {
		return err
	}

	timeout := m.lockForceTimeout
	if timeout > 0 {
		ctx, _ = ldctx.WithTimeout(ctx, timeout)
	}

	ticker := time.NewTicker(m.lockForceInterval)
	defer ticker.Stop()
	for err != nil {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return err

		case <-ticker.C:
			err = m.internalLock(ctx, key, token)
		}
	}

	return nil
}

func (m *Mutex) internalLock(ctx ldctx.Context, key, token string) error {
	cli := m.redis.Client()

	cmd := cli.SetNX(ctx, key, token, m.getExpiration())
	if err := cmd.Err(); err != nil {
		ldctx.LogE(ctx, "redis mutex setnx fail", zap.Error(err), getCallerField(m.redis))
		return lderr.Wrap(err, lderr.ErrCacheRead)
	}

	if ok := cmd.Val(); !ok {
		ldctx.LogW(ctx, "redis mutex has been locked by another goroutine/process", getCallerField(m.redis))
		return lderr.ErrCacheMutexLocked
	}

	return nil
}

func (m *Mutex) Unlock(ctx ldctx.Context) error {
	d := m.unlockDelay
	if d <= 0 {
		return m.unlock(ctx)
	}

	go func() {
		time.Sleep(d)
		m.unlock(ctx)
	}()

	return nil
}

func (m *Mutex) unlock(ctx ldctx.Context) error {
	mc := m.ctx.Load()
	if mc == nil {
		ldctx.LogW(ctx, "redis mutex has not been locked", getCallerField(m.redis))
		return nil
	}

	lockTime := mc.lockTime.Load()
	if lockTime == 0 {
		ldctx.LogW(ctx, "redis mutex has not been locked", getCallerField(m.redis))
		return nil
	}

	// ctx = mc.ctx
	cli := m.redis.Client()
	key := mc.key
	val := mc.token

	ctx = ldctx.WithLogField(ctx, zap.String("key", key), zap.String("token", val))

	if ok := mc.lockTime.CompareAndSwap(lockTime, 0); !ok {
		ldctx.LogW(ctx, "redis mutex has been unlocked by another goroutine", getCallerField(m.redis))
		return nil
	}
	m.ctx.CompareAndSwap(mc, nil)

	ldctx.LogD(ctx, "redis mutex will be unlocked", getCallerField(m.redis))

	// ldctx.TryCancel(mc.ctx)
	mc.cancel()
	if err := m.checkToken(ctx, mc); err != nil {
		return err
	}

	cmd := cli.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		ldctx.LogW(ctx, "redis mutex del fail", zap.Error(err), getCallerField(m.redis))
		return lderr.Wrap(err, lderr.ErrCacheWrite)
	}

	return nil
}

func (m *Mutex) goroutine(mc *mutexContext, lockTime int64) {
	ctx := mc.ctx
	ldctx.LogD(ctx, "redis mutex goroutine start")
	ticker := time.NewTicker(m.heartbeatInterval)

	defer func() {
		ldctx.LogD(ctx, "redis mutex goroutine stop")

		// ldctx.TryCancel(ctx)
		mc.cancel()
		mc.lockTime.CompareAndSwap(lockTime, 0)

		close(mc.events)
		m.ctx.CompareAndSwap(mc, nil)

		ticker.Stop()
	}()

	mc.lastHeartbeat = time.Now()
	for running := true; running; {
		select {
		case now := <-ticker.C:
			running = m.heartbeat(ctx, mc, now)

		case <-ctx.Done():
			c := ldctx.WithLogger(ldctx.Default(), ldctx.GetLogger(ctx))
			m.unlock(c)
			return
		}
	}
}

func (m *Mutex) heartbeat(ctx ldctx.Context, mc *mutexContext, now time.Time) bool {
	switch err := m.checkToken(mc.ctx, mc); err {
	case nil:
		mc.lastHeartbeat = now

	case lderr.ErrCacheMutexNotExists, lderr.ErrCacheMutexNotMatch:
		m.doHeartbeatError(ctx, mc)
		return false

	default:
		return m.checkHeartbeatTime(ctx, mc)
	}

	return true
}

func (m *Mutex) checkToken(ctx ldctx.Context, mc *mutexContext) error {
	cli := m.redis.Client()
	key := mc.key
	val := mc.token
	{
		cmd := cli.Expire(ctx, key, m.getExpiration())
		if err := cmd.Err(); err != nil {
			ldctx.LogE(ctx, "redis mutex expire fail", zap.Error(err))
			return lderr.Wrap(err, lderr.ErrCacheWrite)
		}

		if ok := cmd.Val(); !ok {
			ldctx.LogE(ctx, "redis mutex is not exists")
			return lderr.ErrCacheMutexNotExists
		}
	}

	{
		cmd := cli.Get(ctx, key)
		if err := cmd.Err(); err != nil {
			ldctx.LogE(ctx, "redis mutex get fail", zap.Error(err))
			return lderr.Wrap(err, lderr.ErrCacheRead)
		}

		if val != cmd.Val() {
			ldctx.LogE(ctx, "redis mutex token is not match", zap.String("old", val),
				zap.String("new", cmd.Val()))
			return lderr.ErrCacheMutexNotMatch
		}
	}

	return nil
}

func (m *Mutex) checkHeartbeatTime(ctx ldctx.Context, mc *mutexContext) bool {
	if time.Since(mc.lastHeartbeat) < m.heartbeatTimeout {
		return true
	}

	ldctx.LogW(ctx, "redis mutex fail too much")
	m.doHeartbeatError(ctx, mc)
	return false
}

func (m *Mutex) doHeartbeatError(ctx ldctx.Context, mc *mutexContext) {
	select {
	case <-ctx.Done():
	case mc.events <- MutexEventDeleted:
	}
}

func init() {
	closedMutextEventsChan = make(chan MutexEvent)
	close(closedMutextEventsChan)
}

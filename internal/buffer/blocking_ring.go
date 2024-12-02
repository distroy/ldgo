/*
 * Copyright (C) distroy
 */

package buffer

import (
	"sync"
)

type BlockingRing[T any] struct {
	Buf Ring[T]

	Mu        sync.Mutex
	ReadCond  sync.Cond
	WriteCond sync.Cond
}

func (b *BlockingRing[T]) Init() {
	b.ReadCond.L = &b.Mu
	b.WriteCond.L = &b.Mu
}

func (b *BlockingRing[T]) Close() error { return bRing(b, func() error { return b.Buf.Close() }) }
func (b *BlockingRing[T]) Closed() bool {
	return bRing(b, func() bool {
		ok := b.Buf.Closed()
		if ok {
			b.ReadCond.Broadcast()
			b.WriteCond.Broadcast()
		}
		return ok
	})
}

func (b *BlockingRing[T]) Cap() int  { return bRing(b, func() int { return b.Buf.Cap() }) }
func (b *BlockingRing[T]) Size() int { return bRing(b, func() int { return b.Buf.Size() }) }

func (b *BlockingRing[T]) Write(d []T) (int, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	{
		n, err := b.Buf.Write(d)
		if err != nil {
			return n, err
		}

		if n > 0 {
			b.ReadCond.Signal()
			return n, nil
		}
	}

	b.WriteCond.Wait()

	n, err := b.Buf.Write(d)
	if n > 0 {
		b.ReadCond.Signal()
		if !b.Buf.full {
			b.WriteCond.Signal()
		}
	}
	return n, err
}

func (b *BlockingRing[T]) Read(d []T) (int, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	{
		n, err := b.Buf.Read(d)
		if err != nil {
			return n, err
		}
		if n > 0 {
			b.WriteCond.Signal()
			return n, nil
		}
	}

	b.ReadCond.Wait()

	n, err := b.Buf.Read(d)
	if n > 0 {
		b.WriteCond.Signal()
		if b.Buf.begin != b.Buf.end {
			b.ReadCond.Signal()
		}
	}
	return n, err
}

func (b *BlockingRing[T]) Put(d T) error {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	{
		ok, err := b.Buf.Put(d)
		if err != nil {
			return err
		}
		if ok {
			b.ReadCond.Signal()
			return nil
		}
	}

	b.WriteCond.Wait()
	ok, err := b.Buf.Put(d)
	if ok {
		b.ReadCond.Signal()
		if !b.Buf.full {
			b.WriteCond.Signal()
		}
	}
	return err
}

func (b *BlockingRing[T]) Pop() (T, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	{
		d, ok, err := b.Buf.Pop()
		if err != nil {
			return d, err
		}
		if ok {
			b.WriteCond.Signal()
			return d, nil
		}
	}

	b.ReadCond.Wait()

	d, ok, err := b.Buf.Pop()
	if ok {
		b.WriteCond.Signal()
		if b.Buf.begin != b.Buf.end {
			b.ReadCond.Signal()
		}
	}
	return d, err
}

func bRing[T any, R any](b *BlockingRing[T], f func() R) R {
	b.Mu.Lock()
	r := f()
	b.Mu.Unlock()
	return r
}

/*
 * Copyright (C) distroy
 */

package ldslice

func Split[S ~[]V, V any](s S, n int) []S {
	l := len(s)
	if l == 0 {
		return nil
	}
	if n >= l || n <= 0 {
		return []S{s}
	}
	r := make([]S, 0, (l+n-1)/n)
	for i := 0; i < l; i += n {
		b := i
		e := i + n
		if e > l {
			e = l
		}
		r = append(r, s[b:e])
		// log.Printf("b:%d, e:%d, r:%v", b, e, r)
	}
	return r
}

func SplitFunc[S ~[]V, V any](s S, n int, f func(s S) bool) int {
	l := len(s)
	if l == 0 {
		return 0
	}
	if n >= l || n <= 0 {
		if f != nil {
			f(s)
		}
		return 1
	}
	if f == nil {
		count := (l + n - 1) / n
		return count
	}
	count := 0
	for i := 0; i < l; i += n {
		b := i
		e := i + n
		if e > l {
			e = l
		}
		ss := s[b:e]
		if ok := f(ss); !ok {
			break
		}
		count++
		// log.Printf("b:%d, e:%d, r:%v", b, e, r)
	}
	return count
}

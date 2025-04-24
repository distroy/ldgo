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

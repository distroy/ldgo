/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	type args struct {
		list []int
		size int
	}
	tests := []struct {
		args args
		want [][]int
	}{
		{
			args: args{list: []int{}, size: 2},
			want: nil,
		},
		{
			args: args{list: []int{1, 2}, size: 4},
			want: [][]int{{1, 2}},
		},
		{
			args: args{list: []int{1, 2, 3, 4, 5, 6, 7}, size: 2},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%d", tt.args.list, tt.args.size)
		t.Run(name, func(t *testing.T) {
			got := Split(tt.args.list, tt.args.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`Split(%v, %d) = %v, want=%v`, tt.args.list, tt.args.size, got, tt.want)
				return
			}
		})
	}
}

func TestSplitFunc(t *testing.T) {
	type args struct {
		list []int
		size int
	}
	tests := []struct {
		args args
		want [][]int
	}{
		{
			args: args{list: []int{}, size: 2},
			want: nil,
		},
		{
			args: args{list: []int{1, 2}, size: 4},
			want: [][]int{{1, 2}},
		},
		{
			args: args{list: []int{1, 2, 3, 4, 5, 6, 7}, size: 2},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%d", tt.args.list, tt.args.size)
		t.Run(name, func(t *testing.T) {
			cnt := SplitFunc(tt.args.list, tt.args.size, nil)
			if cnt != len(tt.want) {
				t.Errorf(`SplitFunc(%v, %d, nil) = %v, want=%v`, tt.args.list, tt.args.size, cnt, len(tt.want))
				return
			}
			var got [][]int
			SplitFunc(tt.args.list, tt.args.size, func(ss []int) bool {
				got = append(got, ss)
				return true
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`SplitFunc(%v, %d, func) = %v, want=%v`, tt.args.list, tt.args.size, got, tt.want)
				return
			}
		})
	}
}

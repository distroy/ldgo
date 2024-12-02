/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"fmt"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		str  string
		sep  rune
		want string
	}{
		{
			str:  "",
			sep:  0,
			want: "",
		},
		{
			str:  "HTTP Request",
			sep:  0,
			want: "http_request",
		},
		{
			str:  "  HTTP Request   ",
			sep:  0,
			want: "http_request",
		},
		{
			str:  "HTTP request",
			sep:  0,
			want: "http_request",
		},
		{
			str:  "HTTPRequest",
			sep:  0,
			want: "http_request",
		},
		{
			str:  "HTTPRequest",
			sep:  '-',
			want: "http-request",
		},
		{
			str:  "struct A",
			want: "struct_a",
		},
		{
			str:  "struct A",
			want: "struct_a",
		},
		{
			str:  "structA",
			want: "struct_a",
		},
		{
			str:  "structA",
			sep:  ' ',
			want: "struct a",
		},
	}
	for i, tt := range tests {
		sep := tt.sep
		name := fmt.Sprintf("%d|%c|%s", i, sep, tt.str)
		t.Run(name, func(t *testing.T) {
			if got := ToSnakeCase(tt.str, tt.sep); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		str  string
		seps []rune
		want string
	}{
		{
			str:  "",
			want: "",
		},
		{
			str:  "HTTP Request",
			want: "HTTPRequest",
		},
		{
			str:  "  HTTP Request  ",
			want: "HTTPRequest",
		},
		{
			str:  "struct-a",
			want: "StructA",
		},
		{
			str:  "struct---   a",
			seps: []rune{'_', '-'},
			want: "StructA",
		},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("%d|%s|%v", i, tt.str, tt.seps)
		t.Run(name, func(t *testing.T) {
			if got := ToCamelCase(tt.str, tt.seps); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

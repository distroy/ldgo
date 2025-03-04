/*
 * Copyright (C) distroy
 */

package timeinternal

import (
	"reflect"
	"testing"
	"time"
)

func TestDurationMarshalJSON(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "",
			args:    args{d: 123456789},
			want:    []byte(`"123.456789ms"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DurationMarshalJson(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationMarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DurationMarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDurationUnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "string",
			args:    args{b: []byte(`"2h1m"`)},
			want:    time.Hour*2 + time.Minute*1,
			wantErr: false,
		},
		{
			name:    "int",
			args:    args{b: []byte("12345")},
			want:    12345,
			wantErr: false,
		},
		{
			name:    "int-0",
			args:    args{b: []byte("0")},
			want:    0,
			wantErr: false,
		},
		{
			name:    "hex-x",
			args:    args{b: []byte("0x12345")},
			want:    0x12345,
			wantErr: false,
		},
		{
			name:    "hex-X",
			args:    args{b: []byte("0X12345")},
			want:    0x12345,
			wantErr: false,
		},
		{
			name:    "oct-o",
			args:    args{b: []byte("0o12345")},
			want:    0o12345,
			wantErr: false,
		},
		{
			name:    "oct-O",
			args:    args{b: []byte("0O12345")},
			want:    0o12345,
			wantErr: false,
		},
		{
			name:    "oct",
			args:    args{b: []byte("012345")},
			want:    0o12345,
			wantErr: false,
		},
		{
			name:    "err-str",
			args:    args{b: []byte("abc")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "err-int",
			args:    args{b: []byte("12345f")},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DurationUnmarshalJson(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationUnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DurationUnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

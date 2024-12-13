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
			got, err := DurationMarshalJSON(tt.args.d)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DurationUnmarshalJSON(tt.args.b)
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

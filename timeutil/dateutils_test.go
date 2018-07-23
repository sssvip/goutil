package timeutil

import (
	"reflect"
	"testing"
	"time"
)

func TestFormatTimeStringByTimestamp(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"t", args{1514739661}, "2018-01-01 01:01:01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTimeStringByTimestamp(tt.args.timestamp); got != tt.want {
				t.Errorf("FormatTimeStringByTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDateTimeString(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{time.Unix(1514739661, 0)}, "2018-01-01 01:01:01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateTimeString(tt.args.t); got != tt.want {
				t.Errorf("FormatDateTimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimeByDateStr(t *testing.T) {
	type args struct {
		dateStr string
	}
	ti := time.Unix(1514764800, 0).UTC()
	tests := []struct {
		name string
		args args
		want *time.Time
	}{
		{"", args{"2018-01-01"}, &ti},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTimeByDateStr(tt.args.dateStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimeByDateStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimeByDateTimeStr(t *testing.T) {
	type args struct {
		dateStr string
	}
	ti := time.Unix(1514764800, 0).UTC()
	tests := []struct {
		name string
		args args
		want *time.Time
	}{
		{"", args{"2018-01-01 00:00:00"}, &ti},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTimeByDateTimeStr(tt.args.dateStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimeByDateStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDateString(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{time.Unix(1514739661, 0)}, "2018-01-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateString(tt.args.t); got != tt.want {
				t.Errorf("FormatDateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

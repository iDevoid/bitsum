package util

import (
	"reflect"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "GMT+9",
			args: args{
				d: time.Date(2021, time.April, 12, 10, 55, 23, 0, loc),
			},
			want: 1618221600,
		},
		{
			name: "UTC",
			args: args{
				d: time.Date(2021, time.April, 2, 14, 23, 45, 0, time.UTC),
			},
			want: 1617372000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.d); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashHours(t *testing.T) {
	locTokyo, _ := time.LoadLocation("Asia/Tokyo")
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "UTC only",
			args: args{
				start: time.Date(2021, time.April, 5, 18, 23, 45, 0, time.UTC),
				end:   time.Date(2021, time.April, 6, 1, 23, 45, 0, time.UTC),
			},
			want: []int64{
				1617645600, 1617649200, 1617652800, 1617656400, 1617660000, 1617663600, 1617667200, 1617670800,
			},
		},
		{
			name: "multi timezone",
			args: args{
				start: time.Date(2021, time.April, 2, 18, 23, 45, 0, locJakarta),
				end:   time.Date(2021, time.April, 3, 1, 23, 45, 0, locTokyo),
			},
			want: []int64{
				1617386400, 1617390000, 1617393600, 1617397200, 1617400800, 1617404400, 1617408000, 1617411600,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashHours(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HashHours() = %v, want %v", got, tt.want)
			}
		})
	}
}

package memcache

import (
	"reflect"
	"testing"
	"time"

	"github.com/iDevoid/bitsum/internal/constant/model"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name string
		want CoinsMemCache
	}{
		{
			name: "success",
			want: &coinsCache{
				Data:   make(map[int64][]model.Coin),
				Newest: &model.Coin{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Initialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coinsCache_Increment(t *testing.T) {
	type fields struct {
		Data   map[int64][]model.Coin
		Newest *model.Coin
	}
	type args struct {
		data *model.Coin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		{
			name: "no record before",
			fields: fields{
				Data:   make(map[int64][]model.Coin),
				Newest: &model.Coin{},
			},
			args: args{
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
			want: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
		},
		{
			name: "existing record too",
			fields: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
			args: args{
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   100,
				},
			},
			want: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
					1618221600: {
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   1100,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   1100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &coinsCache{
				Data:   tt.fields.Data,
				Newest: tt.fields.Newest,
			}
			if err := cc.Increment(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("coinsCache.Increment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(cc.Data, tt.want.Data) {
				t.Errorf("coinsCache.Increment() data = %v, want %v", cc.Data, tt.want.Data)
			}
			if !reflect.DeepEqual(cc.Newest, tt.want.Newest) {
				t.Errorf("coinsCache.Increment() newest = %v, want %v", cc.Newest, tt.want.Newest)
			}
		})
	}
}

func Test_coinsCache_Decrement(t *testing.T) {
	type fields struct {
		Data   map[int64][]model.Coin
		Newest *model.Coin
	}
	type args struct {
		data *model.Coin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
			args: args{
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   100,
				},
			},
			want: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
					1618221600: {
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   900,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   900,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &coinsCache{
				Data:   tt.fields.Data,
				Newest: tt.fields.Newest,
			}
			if err := cc.Decrement(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("coinsCache.Decrement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_coinsCache_History(t *testing.T) {
	type fields struct {
		Data   map[int64][]model.Coin
		Newest *model.Coin
	}
	type args struct {
		data *model.FilterDate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Coin
		wantErr bool
	}{
		{
			name: "no record",
			fields: fields{
				Data:   make(map[int64][]model.Coin),
				Newest: &model.Coin{},
			},
			args: args{
				data: &model.FilterDate{
					StartDateTime: time.Date(2021, time.April, 10, 10, 55, 23, 0, time.UTC),
					EndDateTime:   time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
				},
			},
			want: []model.Coin{},
		},
		{
			name: "with record",
			fields: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
					1618221600: {
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   1100,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   1100,
				},
			},
			args: args{
				data: &model.FilterDate{
					StartDateTime: time.Date(2021, time.April, 10, 10, 55, 23, 0, time.UTC),
					EndDateTime:   time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
				},
			},
			want: []model.Coin{
				{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
				{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   1100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &coinsCache{
				Data:   tt.fields.Data,
				Newest: tt.fields.Newest,
			}
			got, err := cc.History(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("coinsCache.History() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("coinsCache.History() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coinsCache_Balance(t *testing.T) {
	type fields struct {
		Data   map[int64][]model.Coin
		Newest *model.Coin
	}
	tests := []struct {
		name    string
		fields  fields
		want    model.Coin
		wantErr bool
	}{
		{
			name: "no record",
			fields: fields{
				Data:   make(map[int64][]model.Coin),
				Newest: &model.Coin{},
			},
			wantErr: true,
		},
		{
			name: "with record",
			fields: fields{
				Data: map[int64][]model.Coin{
					1618135200: {
						{
							DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
							Amount:   1000,
						},
					},
					1618221600: {
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   1100,
						},
					},
				},
				Newest: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   1100,
				},
			},
			want: model.Coin{
				DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
				Amount:   1100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &coinsCache{
				Data:   tt.fields.Data,
				Newest: tt.fields.Newest,
			}
			got, err := cc.Balance()
			if (err != nil) != tt.wantErr {
				t.Errorf("coinsCache.Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("coinsCache.Balance() = %v, want %v", got, tt.want)
			}
		})
	}
}

package coins

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/storage/memcache"
	coins_mock "github.com/iDevoid/bitsum/mocks/coins"
)

func Test_service_HistoryTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		coinsCache memcache.CoinsMemCache
	}
	type args struct {
		ctx  context.Context
		data *model.FilterDate
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCoins []model.Coin
		wantErr   bool
	}{
		{
			name: "bad date",
			args: args{
				ctx: context.TODO(),
				data: &model.FilterDate{
					StartDateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					EndDateTime:   time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
				},
			},
			wantErr: true,
		},
		{
			name: "error",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().History(&model.FilterDate{
						StartDateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
						EndDateTime:   time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					}).Return([]model.Coin{}, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.FilterDate{
					StartDateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					EndDateTime:   time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
				},
			},
			wantCoins: []model.Coin{},
			wantErr:   true,
		},
		{
			name: "success",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().History(&model.FilterDate{
						StartDateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
						EndDateTime:   time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
					}).Return([]model.Coin{
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   10,
						},
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.FilterDate{
					StartDateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					EndDateTime:   time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
				},
			},
			wantCoins: []model.Coin{
				{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				coinsCache: tt.fields.coinsCache,
			}
			gotCoins, err := s.HistoryTransaction(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.HistoryTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCoins, tt.wantCoins) {
				t.Errorf("service.HistoryTransaction() = %v, want %v", gotCoins, tt.wantCoins)
			}
		})
	}
}

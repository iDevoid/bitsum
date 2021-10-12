package coins

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/storage/memcache"
	coins_mock "github.com/iDevoid/bitsum/mocks/coins"
)

func Test_service_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		coinsCache memcache.CoinsMemCache
	}
	type args struct {
		ctx  context.Context
		data *model.Coin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error back date",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   100,
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
			wantErr: true,
		},
		{
			name: "with no record",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{}, errors.New("ERROR"))
					mocked.EXPECT().Increment(&model.Coin{
						DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
						Amount:   1000,
					}).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
			wantErr: true,
		},
		{
			name: "with record success",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 10, 10, 55, 23, 0, time.UTC),
						Amount:   100,
					}, nil)
					mocked.EXPECT().Increment(&model.Coin{
						DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
						Amount:   1000,
					}).Return(nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   1000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				coinsCache: tt.fields.coinsCache,
			}
			if err := s.Receive(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("service.Receive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

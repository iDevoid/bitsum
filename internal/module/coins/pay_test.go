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

func Test_service_Pay(t *testing.T) {
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
			name: "no balance recorded",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{}, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "insufficient balance",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
					Amount:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "error back date",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
					Amount:   10,
				},
			},
			wantErr: true,
		},
		{
			name: "error",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}, nil)
					mocked.EXPECT().Decrement(&model.Coin{
						DateTime: time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
					Amount:   10,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}, nil)
					mocked.EXPECT().Decrement(&model.Coin{
						DateTime: time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}).Return(nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				data: &model.Coin{
					DateTime: time.Date(2021, time.April, 13, 10, 55, 23, 0, time.UTC),
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
			if err := s.Pay(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("service.Pay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

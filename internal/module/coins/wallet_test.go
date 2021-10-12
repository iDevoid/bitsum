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

func Test_service_Balance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		coinsCache memcache.CoinsMemCache
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes model.Coin
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{}, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				coinsCache: func() memcache.CoinsMemCache {
					mocked := coins_mock.NewMockCoinsMemCache(ctrl)
					mocked.EXPECT().Balance().Return(model.Coin{
						DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
						Amount:   1000,
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantRes: model.Coin{
				DateTime: time.Date(2021, time.April, 11, 10, 55, 23, 0, time.UTC),
				Amount:   1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				coinsCache: tt.fields.coinsCache,
			}
			gotRes, err := s.Balance(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("service.Balance() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

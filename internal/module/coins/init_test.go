package coins

import (
	"reflect"
	"testing"

	"github.com/iDevoid/bitsum/internal/storage/memcache"
	coins_mock "github.com/iDevoid/bitsum/mocks/coins"
)

func TestInitialize(t *testing.T) {
	type args struct {
		coinsCache memcache.CoinsMemCache
	}
	tests := []struct {
		name string
		args args
		want Usecase
	}{
		{
			name: "success",
			args: args{
				coinsCache: &coins_mock.MockCoinsMemCache{},
			},
			want: &service{
				coinsCache: &coins_mock.MockCoinsMemCache{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Initialize(tt.args.coinsCache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

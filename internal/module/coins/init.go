package coins

import (
	"context"

	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/storage/memcache"
)

//go:generate mockgen -destination=../../../mocks/coins/usecase_mock.go -package=coins_mock -source=init.go

type Usecase interface {
	// Pay takes amount to decrement the balance and update the latest historical data
	// if date is behind the latest data, it returns error
	// if the balance amount is insufficient, returns error
	Pay(ctx context.Context, data *model.Coin) error
	// Receive takes new data to increment the balance of coins and update the latest date
	// does not allow the request date behind the latest date
	Receive(ctx context.Context, data *model.Coin) error

	HistoryTransaction(ctx context.Context, data *model.FilterDate) (coins []model.Coin, err error)
	// Balance returns the current balence or the latest transactical data
	Balance(ctx context.Context) (res model.Coin, err error)
}

type service struct {
	coinsCache memcache.CoinsMemCache
}

func Initialize(coinsCache memcache.CoinsMemCache) Usecase {
	return &service{
		coinsCache,
	}
}

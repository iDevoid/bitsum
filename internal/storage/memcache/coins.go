package memcache

import (
	"errors"

	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/constant/state"
	"github.com/iDevoid/bitsum/internal/util"
)

//go:generate mockgen -destination=../../../mocks/coins/memcache_mock.go -package=coins_mock -source=coins.go

type CoinsMemCache interface {
	// Increment receives new amount and incrising it as historical data
	Increment(data *model.Coin) error
	// Decrement receives new amount and Decresing it as historical data
	Decrement(data *model.Coin) error
	// History takes range of dates and gets all the data
	History(data *model.FilterDate) (res []model.Coin, err error)
	// Balance returns the latest historical data of coin
	Balance() (model.Coin, error)
}

type coinsCache struct {
	// the index is a hash from date in unix
	Data   map[int64][]model.Coin
	Newest *model.Coin
}

// Initialize the memory caching for coins historical data
func Initialize() CoinsMemCache {
	return &coinsCache{
		Data:   make(map[int64][]model.Coin),
		Newest: &model.Coin{},
	}
}

// Increment receives new amount and incrising it as historical data
func (cc *coinsCache) Increment(data *model.Coin) error {
	hashcode := util.Hash(data.DateTime)

	data.Amount += cc.Newest.Amount
	cc.Data[hashcode] = append(cc.Data[hashcode], *data)

	cc.Newest.Amount = data.Amount
	cc.Newest.DateTime = data.DateTime
	return nil
}

// Decrement receives new amount and Decresing it as historical data
func (cc *coinsCache) Decrement(data *model.Coin) error {
	hashcode := util.Hash(data.DateTime)

	data.Amount = cc.Newest.Amount - data.Amount
	cc.Data[hashcode] = append(cc.Data[hashcode], *data)

	cc.Newest.Amount = data.Amount
	cc.Newest.DateTime = data.DateTime
	return nil
}

// History takes range of dates and gets all the data
func (cc *coinsCache) History(data *model.FilterDate) ([]model.Coin, error) {
	res := make([]model.Coin, 0)
	dates := util.HashHours(data.StartDateTime, data.EndDateTime)
	for _, hash := range dates {
		records, ok := cc.Data[hash]
		if !ok {
			continue
		}
		for _, record := range records {
			curr := record.DateTime.Unix()
			start := data.StartDateTime.Unix()
			end := data.EndDateTime.Unix()
			if curr >= start && start <= end {
				res = append(res, record)
			}
		}
	}
	return res, nil
}

// Balance returns the latest historical data of coin
func (cc *coinsCache) Balance() (model.Coin, error) {
	if len(cc.Data) == 0 {
		return model.Coin{}, errors.New(state.ErrorEmpty)
	}
	return *cc.Newest, nil
}

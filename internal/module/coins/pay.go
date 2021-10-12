package coins

import (
	"context"
	"errors"

	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/constant/state"
)

// Pay takes amount to decrement the balance and update the latest historical data
// if date is behind the latest data, it returns error
// if the balance amount is insufficient, returns error
func (s *service) Pay(ctx context.Context, data *model.Coin) error {
	data.DateTime = data.DateTime.UTC()

	wallet, err := s.coinsCache.Balance()
	if err != nil {
		return err
	}
	if wallet.Amount < data.Amount {
		return errors.New(state.ErrorInsufficient)
	}
	if data.DateTime.Before(wallet.DateTime) {
		return errors.New(state.ErrorBackDate)
	}

	err = s.coinsCache.Decrement(data)
	return err
}

package coins

import (
	"context"
	"errors"

	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/constant/state"
)

// Receive takes new data to increment the balance of coins and update the latest date
// does not allow the request date behind the latest date
func (s *service) Receive(ctx context.Context, data *model.Coin) error {
	data.DateTime = data.DateTime.UTC()

	// no handling error to allow adding record even no history or wallet recorded
	wallet, _ := s.coinsCache.Balance()
	if data.DateTime.Before(wallet.DateTime) {
		return errors.New(state.ErrorBackDate)
	}

	err := s.coinsCache.Increment(data)
	return err
}

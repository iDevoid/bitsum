package coins

import (
	"context"
	"errors"

	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/constant/state"
)

func (s *service) HistoryTransaction(ctx context.Context, data *model.FilterDate) (coins []model.Coin, err error) {
	data.StartDateTime = data.StartDateTime.UTC()
	data.EndDateTime = data.EndDateTime.UTC()

	if data.StartDateTime.After(data.EndDateTime) {
		return nil, errors.New(state.ErrorWrongDate)
	}
	coins, err = s.coinsCache.History(data)
	return
}

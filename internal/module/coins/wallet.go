package coins

import (
	"context"

	"github.com/iDevoid/bitsum/internal/constant/model"
)

// Balance returns the current balence or the latest transactical data
func (s *service) Balance(ctx context.Context) (res model.Coin, err error) {
	res, err = s.coinsCache.Balance()
	return
}

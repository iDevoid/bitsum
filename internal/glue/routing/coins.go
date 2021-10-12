package routing

import (
	"net/http"

	"github.com/iDevoid/bitsum/internal/handler/rest"
	"github.com/iDevoid/bitsum/platform/routers"
)

// CoinsRouting returns the list of routers for domain coins
func CoinsRouting(handler rest.CoinsHandler) []routers.Router {
	return []routers.Router{
		{
			Method:  http.MethodGet,
			Path:    "/wallet",
			Handler: handler.Wallet,
		},

		{
			Method:  http.MethodPost,
			Path:    "/pay",
			Handler: handler.Payment,
		},
		{
			Method:  http.MethodPost,
			Path:    "/receive",
			Handler: handler.Receive,
		},
		{
			Method:  http.MethodPost,
			Path:    "/history",
			Handler: handler.History,
		},
	}
}

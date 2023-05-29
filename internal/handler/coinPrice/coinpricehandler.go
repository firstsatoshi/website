package coinPrice

import (
	"net/http"

	"github.com/fantopia-dev/website/internal/logic/coinPrice"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/response"
)

func CoinPriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coinPrice.NewCoinPriceLogic(r.Context(), svcCtx)
		resp, err := l.CoinPrice()
		response.Response(w, resp, err)
	}
}

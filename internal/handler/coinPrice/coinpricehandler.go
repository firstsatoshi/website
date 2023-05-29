package coinPrice

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/coinPrice"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/response"
)

func CoinPriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coinPrice.NewCoinPriceLogic(r.Context(), svcCtx)
		resp, err := l.CoinPrice()
		response.Response(w, resp, err)
	}
}

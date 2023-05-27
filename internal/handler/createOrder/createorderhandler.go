package createOrder

import (
	"net/http"

	"github.com/fantopia-dev/website/internal/logic/createOrder"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"
	"github.com/fantopia-dev/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := createOrder.NewCreateOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateOrder(&req)
		response.Response(w, resp, err)
	}
}

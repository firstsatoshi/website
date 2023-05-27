package queryOrder

import (
	"net/http"

	"github.com/fantopia-dev/website/internal/logic/queryOrder"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"
	"github.com/fantopia-dev/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := queryOrder.NewQueryOrderLogic(r.Context(), svcCtx)
		resp, err := l.QueryOrder(&req)
		response.Response(w, resp, err)
	}
}

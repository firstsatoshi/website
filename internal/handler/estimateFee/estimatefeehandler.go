package estimateFee

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/estimateFee"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EstimateFeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EstimateFeeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := estimateFee.NewEstimateFeeLogic(r.Context(), svcCtx)
		resp, err := l.EstimateFee(&req)
		response.Response(w, resp, err)
	}
}

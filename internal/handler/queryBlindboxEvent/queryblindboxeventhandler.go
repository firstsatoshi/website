package queryBlindboxEvent

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/queryBlindboxEvent"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryBlindboxEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryBlindboxEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := queryBlindboxEvent.NewQueryBlindboxEventLogic(r.Context(), svcCtx)
		resp, err := l.QueryBlindboxEvent(&req)
		response.Response(w, resp, err)
	}
}

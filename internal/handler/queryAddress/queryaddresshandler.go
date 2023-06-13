package queryAddress

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/queryAddress"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := queryAddress.NewQueryAddressLogic(r.Context(), svcCtx)
		resp, err := l.QueryAddress(&req)
		response.Response(w, resp, err)
	}
}

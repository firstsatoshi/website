package queryBrc20

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/queryBrc20"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryBrc20Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryBrc20Req
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := queryBrc20.NewQueryBrc20Logic(r.Context(), svcCtx)
		resp, err := l.QueryBrc20(&req)
		response.Response(w, resp, err)
	}
}

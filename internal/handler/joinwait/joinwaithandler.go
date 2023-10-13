package joinwait

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/joinwait"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func JoinWaitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinWaitReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := joinwait.NewJoinWaitLogic(r.Context(), svcCtx)
		resp, err := l.JoinWait(&req)
		response.Response(w, resp, err)
	}
}

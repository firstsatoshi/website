package joinwaitlist

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/joinwaitlist"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func JoinWaitListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinWaitListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := joinwaitlist.NewJoinWaitListLogic(r.Context(), svcCtx)
		resp, err := l.JoinWaitList(&req)
		response.Response(w, resp, err)
	}
}

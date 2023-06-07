package checkwhitelist

import (
	"net/http"

	"github.com/firstsatoshi/website/response"
	"github.com/firstsatoshi/website/internal/logic/checkwhitelist"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckWhitelistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckWhitelistReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := checkwhitelist.NewCheckWhitelistLogic(r.Context(), svcCtx)
		resp, err := l.CheckWhitelist(&req)
		response.Response(w, resp, err)
	}
}

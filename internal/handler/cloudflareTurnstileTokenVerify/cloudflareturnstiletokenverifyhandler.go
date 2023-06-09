package cloudflareTurnstileTokenVerify

import (
	"net/http"

	"github.com/firstsatoshi/website/response"
	"github.com/firstsatoshi/website/internal/logic/cloudflareTurnstileTokenVerify"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CloudflareTurnstileTokenVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloudflareTurnstileTokenVerifyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cloudflareTurnstileTokenVerify.NewCloudflareTurnstileTokenVerifyLogic(r.Context(), svcCtx)
		resp, err := l.CloudflareTurnstileTokenVerify(&req)
		response.Response(w, resp, err)
	}
}

package cloudflareTurnstileTokenVerify

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/cloudflareTurnstileTokenVerify"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/response"
)

func CloudflareTurnstileTokenVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := cloudflareTurnstileTokenVerify.NewCloudflareTurnstileTokenVerifyLogic(r.Context(), svcCtx, &r.Form)
		resp, err := l.CloudflareTurnstileTokenVerify()
		response.Response(w, resp, err)
	}
}

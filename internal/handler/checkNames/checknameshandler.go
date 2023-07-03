package checkNames

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/checkNames"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckNamesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckNamesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := checkNames.NewCheckNamesLogic(r.Context(), svcCtx)
		resp, err := l.CheckNames(&req)
		response.Response(w, resp, err)
	}
}

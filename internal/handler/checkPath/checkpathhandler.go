package checkPath

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/checkPath"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckPathHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := checkPath.NewCheckPathLogic(r.Context(), svcCtx)
		resp, err := l.CheckPath(&req)
		response.Response(w, resp, err)
	}
}

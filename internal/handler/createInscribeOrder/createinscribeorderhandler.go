package createInscribeOrder

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/createInscribeOrder"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateInscribeOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateInscribeOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := createInscribeOrder.NewCreateInscribeOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateInscribeOrder(&req)
		response.Response(w, resp, err)
	}
}

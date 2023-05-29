package queryBlindboxEvent

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/queryBlindboxEvent"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/response"
)

func QueryBlindboxEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := queryBlindboxEvent.NewQueryBlindboxEventLogic(r.Context(), svcCtx)
		resp, err := l.QueryBlindboxEvent()
		response.Response(w, resp, err)
	}
}

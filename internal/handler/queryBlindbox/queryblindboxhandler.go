package queryBlindbox

import (
	"net/http"

	"github.com/fantopia-dev/website/internal/logic/queryBlindbox"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/response"
)

func QueryBlindboxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := queryBlindbox.NewQueryBlindboxLogic(r.Context(), svcCtx)
		resp, err := l.QueryBlindbox()
		response.Response(w, resp, err)
	}
}

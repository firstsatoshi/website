package queryGalleryList

import (
	"net/http"

	"github.com/firstsatoshi/website/internal/logic/queryGalleryList"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryGalleryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryGalleryListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := queryGalleryList.NewQueryGalleryListLogic(r.Context(), svcCtx)
		resp, err := l.QueryGalleryList(&req)
		response.Response(w, resp, err)
	}
}

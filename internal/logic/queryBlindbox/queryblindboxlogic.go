package queryBlindbox

import (
	"context"

	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBlindboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryBlindboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBlindboxLogic {
	return &QueryBlindboxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryBlindboxLogic) QueryBlindbox() (resp *types.QueryBlindboxResp, err error) {
	// todo: add your logic here and delete this line

	return
}

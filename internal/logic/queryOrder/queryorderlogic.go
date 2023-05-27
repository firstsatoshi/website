package queryOrder

import (
	"context"

	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOrderLogic {
	return &QueryOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryOrderLogic) QueryOrder(req *types.QueryOrderReq) (resp *types.QueryOrderResp, err error) {
	// todo: add your logic here and delete this line

	return
}

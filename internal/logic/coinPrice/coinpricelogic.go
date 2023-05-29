package coinPrice

import (
	"context"

	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoinPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoinPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoinPriceLogic {
	return &CoinPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoinPriceLogic) CoinPrice() (resp *types.CoinPriceResp, err error) {
	// todo: add your logic here and delete this line

	return
}

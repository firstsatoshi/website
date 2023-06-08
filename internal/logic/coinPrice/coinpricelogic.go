package coinPrice

import (
	"context"
	"fmt"
	"strconv"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"

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

func (l *CoinPriceLogic) CoinPrice() (*types.CoinPriceResp, error) {
	key := "bitcion-price-usd"
	v, err := l.svcCtx.Redis.Get(key)
	if err == nil {
		price, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return &types.CoinPriceResp{
				BtcPriceUsd: float64(int32(price)),
			}, nil
		}
	}

	return nil, fmt.Errorf("coincap api error")
}

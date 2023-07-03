package queryBrc20

import (
	"context"
	"strconv"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBrc20Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryBrc20Logic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBrc20Logic {
	return &QueryBrc20Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryBrc20Logic) QueryBrc20(req *types.QueryBrc20Req) (resp *types.QueryBrc20Resp, err error) {

	brc20Info, err := l.svcCtx.UnisatApiClient.GetBrc20Info(req.Ticker)
	if err != nil || brc20Info == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_BRC20_INFO_ERROR), "get brc20 info error: %v", err.Error())
	}

	if brc20Info.Code != 0 {
		if brc20Info.Code == -1 {
			return &types.QueryBrc20Resp{
				IsExists: false,
				Ticker: req.Ticker,
			}, nil
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_BRC20_INFO_ERROR), "get brc20 info error: %v", brc20Info.Msg)
	}

	limit, err := strconv.ParseInt(brc20Info.Data.Limit, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_BRC20_INFO_ERROR), "parse limit error: %v", err.Error())
	}

	minted, err := strconv.ParseInt(brc20Info.Data.Minted, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_BRC20_INFO_ERROR), "parse minted error: %v", err.Error())
	}

	maxSupply, err := strconv.ParseInt(brc20Info.Data.Max, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_BRC20_INFO_ERROR), "parse max error: %v", err.Error())
	}

	resp = &types.QueryBrc20Resp{
		IsExists:        true,
		Ticker:        brc20Info.Data.Ticker,
		Limit:         limit,
		Minted:        minted,
		Max:           maxSupply,
		Decimal:       brc20Info.Data.Decimal,
		InscriptionId: brc20Info.Data.InscriptionId,
	}

	return
}

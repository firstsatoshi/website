package queryOrder

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"
	"github.com/fantopia-dev/website/xerr"
	"github.com/pkg/errors"

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
	// query exists unpayment order
	queryBuilder := l.svcCtx.TbOrderModel.RowBuilder()

	if len(req.OrderId) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"order_id": req.OrderId,
		}).OrderBy("id ASC")
	} else if len(req.ReceiveAddress) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"receive_address": req.ReceiveAddress,
		}).OrderBy("id ASC")
	} else if len(req.DepositAddress) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"deposit_address": req.DepositAddress,
		}).OrderBy("id ASC")
	} else if len(req.FromAddress) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"pay_from_address": req.FromAddress,
		}).OrderBy("id ASC")
	} else {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "query filter must not be empty")
	}

	orders, err := l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOrders error")
	}

	resp = &types.QueryOrderResp{Orders: make([]types.Order, 0)}
	for _, o := range orders {

		pt := o.PayTime.Time.String()
		if !o.PayTime.Valid {
			pt = ""
		}

		pct := o.PayConfirmedTime.Time.String()
		if !o.PayConfirmedTime.Valid {
			pct = ""
		}
		resp.Orders = append(resp.Orders, types.Order{
			OrderId:          o.OrderId,
			DepositAddress:   o.DepositAddress,
			Total:            int(o.TotalAmountSat),
			ReceiveAddress:   o.ReceiveAddress,
			OrderStatus:      o.OrderStatus,
			PayTime:          pt,
			PayConfirmedTime: pct,
			RevealTxid:       o.RevealTxid.String,
		})
	}

	return
}

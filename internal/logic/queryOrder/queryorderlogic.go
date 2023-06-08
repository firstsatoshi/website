package queryOrder

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
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

func (l *QueryOrderLogic) QueryOrder(req *types.QueryOrderReq) (resp []types.Order, err error) {
	// query exists unpayment order
	queryBuilder := l.svcCtx.TbOrderModel.RowBuilder()

	if len(req.OrderId) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"order_id": req.OrderId,
		}).OrderBy("id DESC")
	} else if len(req.ReceiveAddress) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"receive_address": req.ReceiveAddress,
		}).OrderBy("id DESC")
	} else if len(req.DepositAddress) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"deposit_address": req.DepositAddress,
		}).OrderBy("id DESC")
	} else {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "query filter must not be empty")
	}

	orders, err := l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOrders error")
	}

	resp = make([]types.Order, 0)
	for _, o := range orders {

		payTime := o.PayTime.Time.String()
		if !o.PayTime.Valid {
			payTime = ""
		}

		payConfirmedTime := o.PayConfirmedTime.Time.String()
		if !o.PayConfirmedTime.Valid {
			payConfirmedTime = ""
		}

		payTxid := ""
		if o.PayTxid.Valid {
			payTxid = o.PayTxid.String
		}

		depositAddress := o.DepositAddress
		if o.OrderStatus == "PAYTIMEOUT" || o.OrderStatus == "PAYPENDING" || o.OrderStatus == "ALLSUCCESS" || o.OrderStatus == "MINTING" {
			// don't show deposit for finished status order
			depositAddress = ""
		}

		nftDetails := make([]types.NftDetail, 0)
		if o.OrderStatus == "MINTING" || o.OrderStatus == "ALLSUCCESS" {
			builder := l.svcCtx.TbLockOrderBlindboxModel.RowBuilder().Where(squirrel.Eq{
				"order_id": o.OrderId,
			})
			lks, err := l.svcCtx.TbLockOrderBlindboxModel.FindAll(l.ctx, builder, "")
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindAll error:%v", err.Error())
			}
			// logx.Infof("lks === %v", len(lks))

			for _, lk := range lks {
				b, err := l.svcCtx.TbBlindboxModel.FindOne(l.ctx, lk.BlindboxId)
				if err != nil {
					return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOne error:%v", err.Error())
				}
				if b.RevealTxid.Valid {
					nftDetails = append(nftDetails, types.NftDetail{
						Txid:        b.RevealTxid.String,
						TxConfirmed: b.Status == "MINT",
						Name:        b.Name,
						Category:    b.Category,
						Description: b.Description,
						ImageUrl:    b.ImgUrl,
					})
				}
			}
		}

		resp = append(resp, types.Order{
			OrderId:          o.OrderId,
			EventId:          int(o.EventId),
			DepositAddress:   depositAddress,
			Total:            int(o.TotalAmountSat),
			ReceiveAddress:   o.ReceiveAddress,
			OrderStatus:      o.OrderStatus,
			PayTime:          payTime,
			PayConfirmedTime: payConfirmedTime,
			NftDetails:       nftDetails,
			PayTxid:          payTxid,
			CreateTime:       o.CreateTime.Format("2006-01-02 15:04:05 +0800 CST"),
		})
	}

	return
}

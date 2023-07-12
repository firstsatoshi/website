package queryOrder

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
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

	if strings.HasPrefix(req.OrderId, "I") {
		logx.Infof("xxxxxxxxxxxxxxxxxxxxxxxxxxx")
		// query exists unpayment order
		queryBuilder := l.svcCtx.TbInscribeOrderModel.RowBuilder()

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

		// inscribe order
		orders, err := l.svcCtx.TbInscribeOrderModel.FindOrders(l.ctx, queryBuilder)
		if err != nil {
			logx.Errorf("FindOrders error: %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOrders error")
		}

		resp = make([]types.Order, 0)
		for _, o := range orders {

			payTime := int64(0)
			if o.PayTime.Valid {
				payTime = o.PayTime.Time.Unix()
			}

			payConfirmedTime := int64(0)
			if o.PayConfirmedTime.Valid {
				payConfirmedTime = o.PayConfirmedTime.Time.Unix()
			}

			payTxid := ""
			if o.PayTxid.Valid {
				payTxid = o.PayTxid.String
			}

			depositAddress := o.DepositAddress

			if req.OrderId == "" {
				// !for safety, only display deposit address , when query by orderId,
				depositAddress = ""
			} else if o.OrderStatus == "PAYTIMEOUT" || o.OrderStatus == "PAYPENDING" || o.OrderStatus == "PAYSUCCESS" ||
				o.OrderStatus == "ALLSUCCESS" || o.OrderStatus == "MINTING" {
				// !for safety , don't show deposit address for paypending/finished status order
				depositAddress = ""
			} else if o.OrderStatus == "NOTPAID" {
			}

			nftDetails := make([]types.NftDetail, 0)
			if o.OrderStatus == "MINTING" || o.OrderStatus == "ALLSUCCESS" {
				builder := l.svcCtx.TbInscribeDataModel.RowBuilder().Where(squirrel.Eq{
					"order_id": o.OrderId,
				})
				datas, err := l.svcCtx.TbInscribeDataModel.FindInscribeDatas(l.ctx, builder)
				if err != nil {
					return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindAll error:%v", err.Error())
				}
				// logx.Infof("lks === %v", len(lks))

				for _, b := range datas {
					if b.RevealTxid.Valid {
						nftDetails = append(nftDetails, types.NftDetail{
							Txid:        b.RevealTxid.String,
							TxConfirmed: b.Status == "MINT",
							Name:        b.FileName,
							Category:    b.ContentType,
							// Description: b.ContentType,
							// ImageUrl:    "",
							FileName:    b.FileName,
							ContentType: b.ContentType,
							Inscription: b.Data,
						})
					}
				}
			}

			resp = append(resp, types.Order{
				OrderId:          o.OrderId,
				EventId:          0,
				DepositAddress:   depositAddress,
				Count:            int(o.Count),
				Total:            int(o.TotalAmountSat),
				FeeRate:          int(o.FeeRate),
				ReceiveAddress:   o.ReceiveAddress,
				OrderStatus:      o.OrderStatus,
				PayTime:          payTime,
				PayConfirmedTime: payConfirmedTime,
				NftDetails:       nftDetails,
				PayTxid:          payTxid,
				CreateTime:       o.CreateTime.Unix(),
				// CreateTime:       o.CreateTime.Format("2006-01-02 15:04:05 +0800 CST"),
			})

		}

	} else {
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

			payTime := int64(0)
			if o.PayTime.Valid {
				payTime = o.PayTime.Time.Unix()
			}

			payConfirmedTime := int64(0)
			if o.PayConfirmedTime.Valid {
				payConfirmedTime = o.PayConfirmedTime.Time.Unix()
			}

			payTxid := ""
			if o.PayTxid.Valid {
				payTxid = o.PayTxid.String
			}

			depositAddress := o.DepositAddress

			if req.OrderId == "" {
				// !for safety, only display deposit address , when query by orderId,
				depositAddress = ""
			} else if o.OrderStatus == "PAYTIMEOUT" || o.OrderStatus == "PAYPENDING" || o.OrderStatus == "PAYSUCCESS" ||
				o.OrderStatus == "ALLSUCCESS" || o.OrderStatus == "MINTING" {
				// !for safety , don't show deposit address for paypending/finished status order
				depositAddress = ""
			} else if o.OrderStatus == "NOTPAID" {
				// FIX: for NOTPAID,  the frontend must queryorder frequently(1s) to get the latest info, to avoid available being NOT ENOUGH!
				event, err := l.svcCtx.TbBlindboxEventModel.FindOne(l.ctx, int64(o.EventId))
				if err != nil {
					if err == model.ErrNotFound {
						return nil, errors.Wrapf(xerr.NewErrCode(xerr.EVENT_NOT_EXISTS_ERROR), "event id does not exists %v", o.EventId)
					} else {
						logx.Errorf("FindOne error:%v", err.Error())
						return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOne error: %v", err.Error())
					}
				}

				// check safety avail
				// query PAYPENDING order as pendingOrders, avail = event.avail - len(pendingOrders)
				queryBuilder := l.svcCtx.TbOrderModel.RowBuilder().Where(squirrel.Eq{
					"event_id":     o.EventId,
					"order_status": "PAYPENDING",
				}).OrderBy("id DESC")
				payPendingOrders, err := l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)
				if err != nil {
					logx.Errorf("FindOrders error: %v", err.Error())
					return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
				}
				orderCounter := 0
				for _, o := range payPendingOrders {
					if o.EventId == event.Id {
						orderCounter += 1
					}
				}
				safeAvail := event.Avail - int64(orderCounter)
				if safeAvail < 0 {
					safeAvail = 0
				}

				// !NOTE: When safeAvail is NOT enough, for safety, DO NOT display depositaddress any more
				if safeAvail < o.Count {
					depositAddress = ""
				}
			} else {
				// Display deposit address
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
							// ImageUrl:    b.ImgUrl,

							Inscription: b.Data,
						})
					}
				}
			}

			resp = append(resp, types.Order{
				OrderId:          o.OrderId,
				EventId:          int(o.EventId),
				DepositAddress:   depositAddress,
				Count:            int(o.Count),
				Total:            int(o.TotalAmountSat),
				FeeRate:          int(o.FeeRate),
				ReceiveAddress:   o.ReceiveAddress,
				OrderStatus:      o.OrderStatus,
				PayTime:          payTime,
				PayConfirmedTime: payConfirmedTime,
				NftDetails:       nftDetails,
				PayTxid:          payTxid,
				CreateTime:       o.CreateTime.Unix(),
				// CreateTime:       o.CreateTime.Format("2006-01-02 15:04:05 +0800 CST"),
			})
		}
	}

	return
}

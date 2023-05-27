package createOrder

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/fantopia-dev/website/common/uniqueid"
	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"
	"github.com/fantopia-dev/website/model"
	"github.com/fantopia-dev/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	// check receiveAddress is valid P2TR address
	_, err = btcutil.DecodeAddress(req.ReceiveAddress, &chaincfg.MainNetParams)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCP2TRADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
	}

	// TODO get mempool recommanded feerate
	// https://mempool.space/api/v1/fees/recommended
	if req.FeeRate < 10 || req.FeeRate > 200 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.FEERATE_TOO_SMALL_ERROR), "feeRate too small %v", req.FeeRate)
	}

	// query exists unpayment order
	// queryBuilder := l.svcCtx.TbOrderModel.RowBuilder().Where(squirrel.Eq{
	// 	"receive_address": req.ReceiveAddress,
	// 	"order_status": "NOPAY",
	// }).OrderBy("id ASC")
	// l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)

	orderId := uniqueid.GenSn("BX")

	// TODO
	depositAddress := "bc1p2yzcv24v9tpw6ffhkqcq994y8p4ps2xfv65wx7nsmg4meuvzd0fqyesxg7"

	createTime := time.Now()
	ord := model.TbOrder{
		OrderId:         orderId,
		DepositAddress:  depositAddress,
		ReceiveAddress:  req.ReceiveAddress,
		InscriptionData: "TODO", // TODO
		FeeRate:         int64(req.FeeRate),
		TxfeeAmountSat:  123456,  // TODO
		ServiceFeeSat:   123456,  // TODO
		PriceSat:        123456,  // TODO
		TotalAmountSat:  1123456, // TODO
		OrderStatus:     "NOPAY",
		CreateTime:      createTime,
	}
	_, err = l.svcCtx.TbOrderModel.Insert(l.ctx, &ord)
	if err != nil {
		logx.Errorf("insert error:%v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert error: %v", err.Error())
	}

	resp = &types.CreateOrderResp{
		OrderId:        ord.OrderId,
		DepositAddress: ord.DepositAddress,
		ReceiveAddress: ord.ReceiveAddress,
		FeeRate:        int(ord.FeeRate),
		Bytes:          12345, // TODO
		InscribeFee:    int(ord.TxfeeAmountSat),
		ServiceFee:     int(ord.ServiceFeeSat),
		Price:          int(ord.PriceSat),
		Total:          int(ord.TotalAmountSat),
		// CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		CreateTime:     createTime.Format("2006-01-02 15:04:05 +0800 CST"),
	}

	return
}

package createOrder

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
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

	// check count
	if req.Count <= 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "count is invalid %v", req.Count)
	}
	if req.Count > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.COUNT_EXCEED_PER_ORDER_LIMIT_ERROR), "count is too large %v", req.Count)
	}

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

	event, err := l.svcCtx.TbBlindboxEventModel.FindOne(l.ctx, int64(req.EventId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.EVENT_NOT_EXISTS_ERROR), "event id does not exists %v", req.EventId)
		} else {
			logx.Errorf("FindOne error:%v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOne error: %v", err.Error())
		}
	}

	// check available count
	if event.Avail < int64(req.Count) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.AVAILABLE_COUNT_IS_NOT_ENOUGH),
			"avail count %v is not enough %v", event.Avail, req.Count)
	}

	// query exists unpayment order
	// queryBuilder := l.svcCtx.TbOrderModel.RowBuilder().Where(squirrel.Eq{
	// 	"receive_address": req.ReceiveAddress,
	// 	"order_status": "NOPAY",
	// }).OrderBy("id ASC")
	// l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)

	// random generate account_index and address_index
	rand.Seed(time.Now().UnixNano())

	accountIndex := rand.Uint32()
	addressIndex := rand.Uint32()
	for {
		_, e := l.svcCtx.TbAddressModel.FindOneByCoinTypeAccountIndexAddressIndex(l.ctx, "BTC", int64(accountIndex), int64(addressIndex))
		if e == model.ErrNotFound {
			break
		}

		// if already exists continue generate random index
		accountIndex = rand.Uint32()
		addressIndex = rand.Uint32()
	}

	_, depositAddress, err := l.svcCtx.KeyManager.GetWifKeyAndAddresss(
		accountIndex,
		addressIndex,
		chaincfg.MainNetParams)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "generate address error %v", err.Error())
	}

	addresInsertResult, err := l.svcCtx.TbAddressModel.Insert(l.ctx, &model.TbAddress{
		Address:      depositAddress,
		CoinType:     "BTC",
		AccountIndex: int64(accountIndex),
		AddressIndex: int64(addressIndex),
	})
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert address error %v", err.Error())
	}
	addressId, err := addresInsertResult.LastInsertId()
	if err != nil {
		addressId = 0
	}

	prefix := "BX" + req.ReceiveAddress[4:8] + req.ReceiveAddress[len(req.ReceiveAddress)-4:] +
		depositAddress[4:8] + depositAddress[len(depositAddress)-4:] +
		fmt.Sprintf("%02d", req.Count) + fmt.Sprintf("%02d", req.FeeRate)
	prefix = strings.ToUpper(prefix)
	orderId := uniqueid.GenSn(prefix)

	createTime := time.Now()
	ord := model.TbOrder{
		OrderId:         orderId,
		EventId:         int64(req.EventId),
		Count:           int64(req.Count),
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
		// TODO: use transaction ?
		l.svcCtx.TbAddressModel.Delete(l.ctx, addressId)

		logx.Errorf("insert error:%v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert error: %v", err.Error())
	}

	resp = &types.CreateOrderResp{
		OrderId:        ord.OrderId,
		EventId:        int(ord.EventId),
		Count:          int(ord.Count),
		DepositAddress: ord.DepositAddress,
		ReceiveAddress: ord.ReceiveAddress,
		FeeRate:        int(ord.FeeRate),
		Bytes:          12345, // TODO
		InscribeFee:    int(ord.TxfeeAmountSat),
		ServiceFee:     int(ord.ServiceFeeSat),
		Price:          int(ord.PriceSat),
		Total:          int(ord.TotalAmountSat),
		// CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		CreateTime: createTime.Format("2006-01-02 15:04:05 +0800 CST"),
	}

	return
}

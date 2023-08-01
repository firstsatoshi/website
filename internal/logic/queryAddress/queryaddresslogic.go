package queryAddress

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/firstsatoshi/website/common/globalvar"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAddressLogic {
	return &QueryAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAddressLogic) QueryAddress(req *types.QueryAddressReq) (*types.QueryAddressResp, error) {

	resp := types.QueryAddressResp{}

	// check receiveAddress
	_, err := btcutil.DecodeAddress(req.ReceiveAddress, l.svcCtx.ChainCfg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
	}

	// check whitelist
	whiteAddress, err := l.svcCtx.TbWaitlistModel.FindOneByEventIdBtcAddress(l.ctx, int64(req.EventId), req.ReceiveAddress)
	if err != nil {
		if err == model.ErrNotFound {
			resp.IsWhitelist = false
		} else {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByBtcAddress error: %v", err.Error())
		}
	} else {
		resp.IsWhitelist = true
	}

	// get address current orders
	// each address can't mint over mint limit
	tmpBuilder := l.svcCtx.TbOrderModel.SumBuilder("`count`").Where(
		"receive_address=?", req.ReceiveAddress,
	)
	tmpBuilder = tmpBuilder.Where("(order_status=? OR order_status=? OR order_status=? OR order_status=? OR order_status=?)",
		"NOTPAID", "PAYPENDING", "PAYSUCCESS", "MINTING", "ALLSUCCESS")

	total, err := l.svcCtx.TbOrderModel.FindSum(l.ctx, tmpBuilder)
	if err != nil {
		logx.Errorf("FindSum error:%v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindCount error: %v", err.Error())
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

	resp.CurrentOrdersTotal = int(total)
	resp.EventId = int(event.Id)
	resp.EventMintLimit = int(event.MintLimit)

	// bitcoinfish
	if  event.EventEndpoint == "bitcoinfish" {
		resp.EventMintLimit = int(whiteAddress.MintLimit)

		// current mint fishes count
		if true {
			// each address can't mint over mint limit
			tmpBuilder := l.svcCtx.TbInscribeOrderModel.SumBuilder("`count`").Where(
				"receive_address=?", req.ReceiveAddress,
			).Where(
				"version=?", globalvar.BITCOIN_FISH_MAGIC_NUMBER,
			)
			tmpBuilder = tmpBuilder.Where("(order_status=? OR order_status=? OR order_status=? OR order_status=? OR order_status=?)",
				"NOTPAID", "PAYPENDING", "PAYSUCCESS", "MINTING", "ALLSUCCESS")
			mintCountSum, err := l.svcCtx.TbInscribeOrderModel.FindSum(l.ctx, tmpBuilder)
			if err != nil {
				logx.Errorf("FindSum error:%v", err.Error())
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindCount error: %v", err.Error())
			}

			logx.Infof("=== mintCountSum: %v", mintCountSum)

			resp.CurrentOrdersTotal = int(mintCountSum)
		}
	}

	return &resp, nil
}

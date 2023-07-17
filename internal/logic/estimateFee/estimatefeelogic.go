package estimateFee

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/firstsatoshi/website/common/ordinals"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"

	"github.com/zeromicro/go-zero/core/logx"
)

type EstimateFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEstimateFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimateFeeLogic {
	return &EstimateFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// EstimateFee to estimate fee of inscribe order
func (l *EstimateFeeLogic) EstimateFee(req *types.EstimateFeeReq) (resp *types.EstimateFeeResp, err error) {

	_, err = btcutil.DecodeAddress(req.ReceiveAddress, l.svcCtx.ChainCfg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
	}

	utxoSat := 546
	inscriptionRequests := make([]ordinals.InscriptionData, 0)
	for _, v := range req.FileUploads {
		// eg: dataURL, err := dataurl.DecodeString(`data:text/plain;charset=utf-8;base64,aGV5YQ==`)
		dataURL, err := dataurl.DecodeString(v.DataUrl)
		if err != nil {
			logx.Errorf("parse dataUrl error error: %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "parse dataUrl error: %v", err.Error())
		}

		// tmpAddr := "bc1p0ftnthhe6gsthnhd6mswg96aukn888tzrqldz0wkmeeewpr4lkus0vqflq"
		// if l.svcCtx.ChainCfg.Name == "testnet3" {
		// 	tmpAddr = "tb1p0ftnthhe6gsthnhd6mswg96aukn888tzrqldz0wkmeeewpr4lkuscykx90"
		// }
		// check receiveAddress

		if err != nil {
			logx.Errorf("insert inscribedata error %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert inscribedata error %v", err.Error())
		}

		inscriptionRequests = append(inscriptionRequests, ordinals.InscriptionData{
			ContentType: dataURL.ContentType(),
			Body:        dataURL.Data,
			Destination: req.ReceiveAddress,
		})
	}
	totalFee, _, err := ordinals.EstimateFee(l.svcCtx.ChainCfg, req.FeeRate, true, inscriptionRequests, int64(utxoSat))
	if err != nil {
		logx.Errorf("EstimateFee error %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "EstimateFee error %v", err.Error())
	}

	resp = &types.EstimateFeeResp{
		TotalFee: totalFee,
	}

	return
}

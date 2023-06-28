package createInscribeOrder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/globalvar"
	"github.com/firstsatoshi/website/common/turnslite"
	"github.com/firstsatoshi/website/common/uniqueid"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInscribeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInscribeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInscribeOrderLogic {
	return &CreateInscribeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// https://bitcoinops.org/en/tools/calc-size/
// https://blockchain-academy.hs-mittweida.de/2023/02/calculation-of-bitcoin-transaction-fees-explained/
// for p2tr:
//
//	Transaction size = Overhead +  57.5 * inputsNumber + 43 * outputsNumber
//	    eg: 1 input 1 output: 10.5 + 57.5 * 1 + 43 * 1 = 111 bytes
func calcFee(utxoSat, imgBytes, count, feeRate float64) int64 {
	averageFileSize := float64(imgBytes)

	utxoOutputValue := float64(utxoSat)
	commitTxSize := float64(68 + (43 + 1))
	commitTxSize += 64
	revealTxSize := 10.5 + (57.5 + 43.0)
	revealTxSize += 64
	feeSats := math.Ceil((averageFileSize/4 + commitTxSize + revealTxSize) * feeRate)
	feeSats = 1000 * math.Ceil(feeSats/1000)

	feeSats *= count

	// base fee
	baseService := 1000 * math.Ceil(feeSats*0.1/1000)
	feeSats += baseService

	total := feeSats + utxoOutputValue*count
	return int64(total)
}

func (l *CreateInscribeOrderLogic) CreateInscribeOrder(req *types.CreateInscribeOrderReq) (resp *types.CreateInscribeOrderResp, err error) {
	// verify cloudflare Turnstile token
	ok, err := turnslite.VeifyToken(l.ctx, req.Token, l.svcCtx.Redis)
	if !ok || err != nil {
		if l.svcCtx.ChainCfg.Net == wire.TestNet3 {
			logx.Infof("============testnet skip token verify==============")
		} else {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_TOKEN_ERROR), "token: %v not exists", req.Token)
		}
	}

	// rate limit
	s := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(s[:])
	code, err := l.svcCtx.PeriodLimit.TakeCtx(l.ctx, "createinscribeorderapiperiodlimit:"+tokenHash)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "PeriodLimit.TakeCtx error: %v", err.Error())
	}
	if code != limit.Allowed {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TOO_MANY_REQUEST_ERROR), "rate limit error: %v", req.ReceiveAddress)
	}

	// check count
	if req.Count <= 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "count is invalid %v", req.Count)
	}
	if req.Count > 50 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.COUNT_EXCEED_PER_ORDER_LIMIT_ERROR), "count is too large %v", req.Count)
	}

	// check receiveAddress is valid P2TR address
	_, err = btcutil.DecodeAddress(req.ReceiveAddress, l.svcCtx.ChainCfg)
	if err != nil || len(req.ReceiveAddress) != 62 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCP2TRADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
	}
	if l.svcCtx.ChainCfg.Net == wire.MainNet {
		if !strings.HasPrefix(req.ReceiveAddress, "bc1p") {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCP2TRADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
		}
	} else {
		// testnet3
		if !strings.HasPrefix(req.ReceiveAddress, "tb1p") {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCP2TRADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
		}
	}

	// TODO get mempool recommanded feerate
	// https://mempool.space/api/v1/fees/recommended
	if l.svcCtx.ChainCfg.Name == chaincfg.MainNetParams.Name && (req.FeeRate < 5 || req.FeeRate > 300) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.FEERATE_TOO_SMALL_ERROR), "feeRate too small %v", req.FeeRate)
	}

	// random generate account_index and address_index
	time.Sleep(time.Microsecond * 1)
	rand.Seed(time.Now().UnixNano() + int64(req.Count) + int64(req.FeeRate) + int64(req.ReceiveAddress[10]) + int64(req.ReceiveAddress[17]))

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
		addressIndex)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "generate address error %v", err.Error())
	}

	addresInsertResult, err := l.svcCtx.TbAddressModel.Insert(l.ctx, &model.TbAddress{
		Address:      depositAddress,
		CoinType:     globalvar.BTC,
		BussinesType: globalvar.BussinesTypeInscribe,
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

	utxoSat := 546
	prefix := "I"
	if l.svcCtx.ChainCfg.Net == wire.TestNet3 {
		// for Testnet
		prefix += "T"
		// utxoSat = 1000
	} else {
		// for Mainnet
		prefix += "M"
		// utxoSat = 10000
	}

	prefix += req.ReceiveAddress[4:8] + req.ReceiveAddress[len(req.ReceiveAddress)-4:] +
		depositAddress[4:8] + depositAddress[len(depositAddress)-4:] +
		fmt.Sprintf("%02d", req.Count) + fmt.Sprintf("%02d", req.FeeRate)
	prefix = strings.ToUpper(prefix)
	orderId := uniqueid.GenSn(prefix)

	// dataURL, err := dataurl.DecodeString(`data:text/plain;charset=utf-8;base64,aGV5YQ==`)
	dataSize := 0
	for _, v := range req.FileUploads {
		dataURL, err := dataurl.DecodeString(v.DataUrl)
		if err != nil {
			logx.Errorf("parse dataUrl error error: %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "parse dataUrl error: %v", err.Error())
		}

		dataSize += len(dataURL.Data)

		// insert inscribe data into db
		_, err = l.svcCtx.TbInscribeDataModel.Insert(l.ctx, &model.TbInscribeData{
			OrderId:     orderId,
			Data:        string(dataURL.Data), // is it ok?
			ContentType: dataURL.ContentType(),
			FileName:    v.FileName,
		})
		if err != nil {
			logx.Errorf("insert inscribedata error %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert inscribedata error %v", err.Error())
		}
	}

	totalFee := calcFee(float64(utxoSat), float64(dataSize), float64(req.Count), float64(req.FeeRate))
	logx.Infof("==========totalFee : %v", totalFee)
	logx.Infof("==========net name: %v", l.svcCtx.ChainCfg.Name)

	createTime := time.Now()
	ord := model.TbInscribeOrder{
		OrderId:        orderId,
		Count:          int64(req.Count),
		DepositAddress: depositAddress,
		ReceiveAddress: req.ReceiveAddress,
		DataBytes:      int64(dataSize),
		FeeRate:        int64(req.FeeRate),
		TxfeeAmountSat: totalFee,
		ServiceFeeSat:  0,
		TotalAmountSat: totalFee,
		OrderStatus:    "NOTPAID",
		Version:        0,
		CreateTime:     createTime,
	}
	_, err = l.svcCtx.TbInscribeOrderModel.Insert(l.ctx, &ord)
	if err != nil {
		// TODO: use transaction ?
		l.svcCtx.TbAddressModel.Delete(l.ctx, addressId)

		logx.Errorf("insert error:%v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert error: %v", err.Error())
	}

	// update bloomFilter for deposit
	if err = l.svcCtx.DepositBloomFilter.Add([]byte(depositAddress)); err != nil {
		time.Sleep(time.Millisecond * 17)
		l.svcCtx.DepositBloomFilter.Add([]byte(depositAddress))
	}

	resp = &types.CreateInscribeOrderResp{
		OrderId:        ord.OrderId,
		Count:          int(ord.Count),
		DepositAddress: ord.DepositAddress,
		ReceiveAddress: ord.ReceiveAddress,
		FeeRate:        int(ord.FeeRate),
		Bytes:          dataSize,
		InscribeFee:    int(ord.TxfeeAmountSat),
		ServiceFee:     int(ord.ServiceFeeSat),
		Total:          int(ord.TotalAmountSat),
		CreateTime:     createTime.Format("2006-01-02 15:04:05 +0800 CST"),
	}

	return
}

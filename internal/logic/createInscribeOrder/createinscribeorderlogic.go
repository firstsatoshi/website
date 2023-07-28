package createInscribeOrder

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/firstsatoshi/website/common/ordinals"
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

	count := len(req.FileUploads)

	// check DataUrl https://www.rfc-editor.org/rfc/rfc2397
	BITCOIN_FISH_MAGIC_NUMBER := 137   // tb_inscribe_order   version   for bitcoinfish
	versionNumber := 0

	totalBytesSize := 0
	bitfishTotalPrice := int64(0)
	bitfishFilesCount := 0
	bitfishMintLimit := 0
	for _, v := range req.FileUploads {
		if !strings.HasPrefix(v.DataUrl, "data:") || !strings.Contains(v.DataUrl, ";base64,") {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "invalid dataUrl %v", v.DataUrl)
		}

		// as filename is varchar(100)
		if len(v.FileName) > 90 || len(v.FileName) == 0 {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "filename too long or empty %v", v.FileName)
		}

		// bitcoinfish
		prefix := "bitcoinfish_"
		suffix := ".html"
		endpoint := "bitcoinfish"
		if strings.HasPrefix(v.FileName, prefix) && strings.HasSuffix(v.FileName, suffix) {
			// parse path
			b64Path, _ := strings.CutPrefix(v.FileName, prefix)
			b64Path, _ = strings.CutSuffix(b64Path, suffix)
			mergePath, err := base64.StdEncoding.DecodeString(b64Path)
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "invalid bitcoinfish filename: %v", v.FileName)
			}

			// checkpath
			_, err = l.svcCtx.TbBitfishMergePathModel.FindOneByMergePath(l.ctx, string(mergePath))
			if err != model.ErrNotFound {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByMergePath error: %v", err.Error())
			}

			// checkwhitelist
			event, err := l.svcCtx.TbBlindboxEventModel.FindOneByEventEndpoint(l.ctx, endpoint)
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByEventEndpoint error: %v", err.Error())
			}
			if event.OnlyWhitelist > 0 {
				whitelist, err := l.svcCtx.TbWaitlistModel.FindOneByEventIdBtcAddress(l.ctx, event.Id, req.ReceiveAddress)
				if err != nil {
					if err == model.ErrNotFound {
						return nil, errors.Wrapf(xerr.NewErrCode(xerr.ONLY_WHITELIST_ERROR), "FindOneByEventIdBtcAddress error: %v", err.Error())
					}
					return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByEventIdBtcAddress error: %v", err.Error())
				}

				// bitfish: how many fishes you stake, how many fishes you could mint
				bitfishMintLimit = int(whitelist.MintLimit)
			}
			bitfishFilesCount += 1

			// bitfish price
			bitfishTotalPrice += event.PriceSats

			// magic number for flag bitcoinfish order
			versionNumber = BITCOIN_FISH_MAGIC_NUMBER
		}

		totalBytesSize += len(v.DataUrl)
	}

	// bitcoinfish
	if true {
		// each address can't mint over mint limit
		tmpBuilder := l.svcCtx.TbInscribeOrderModel.SumBuilder("`count`").Where(
			"receive_address=?", req.ReceiveAddress,
		).Where(
			"version=?", BITCOIN_FISH_MAGIC_NUMBER,
		)
		tmpBuilder = tmpBuilder.Where("(order_status=? OR order_status=? OR order_status=? OR order_status=? OR order_status=?)",
			"NOTPAID", "PAYPENDING", "PAYSUCCESS", "MINTING", "ALLSUCCESS")
		mintCountSum, err := l.svcCtx.TbInscribeOrderModel.FindSum(l.ctx, tmpBuilder)
		if err != nil {
			logx.Errorf("FindSum error:%v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindCount error: %v", err.Error())
		}

		logx.Infof("=== bitfishMintLimit: %v,  bitfishFilesCount: %v, mintCountSum: %v", bitfishMintLimit, bitfishFilesCount, mintCountSum)
		if bitfishFilesCount+int(mintCountSum) > bitfishMintLimit {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.EXCEED_MINT_LIMIT_ERROR), "address exceed mint limit %v", bitfishMintLimit)
		}
	}

	// max limit size  1MB
	if totalBytesSize > 1_000_001 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "total bytes size too large %v", totalBytesSize)
	}

	// check count
	if count <= 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "empty inscription %v", count)
	}
	if count > 100 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.COUNT_EXCEED_PER_ORDER_LIMIT_ERROR), "count is too large %v", count)
	}

	// check receiveAddress
	_, err = btcutil.DecodeAddress(req.ReceiveAddress, l.svcCtx.ChainCfg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCADDRESS_ERROR), "invalid receive address %v", req.ReceiveAddress)
	}

	// TODO get mempool recommanded feerate
	// https://mempool.space/api/v1/fees/recommended
	if l.svcCtx.ChainCfg.Name == chaincfg.MainNetParams.Name && (req.FeeRate < 3 || req.FeeRate > 300) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.FEERATE_TOO_SMALL_ERROR), "feeRate too small %v", req.FeeRate)
	}

	// random generate account_index and address_index
	time.Sleep(time.Microsecond * 1)
	rand.Seed(time.Now().UnixNano() + int64(count) + int64(req.FeeRate) + int64(req.ReceiveAddress[10]) + int64(req.ReceiveAddress[17]))

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
		fmt.Sprintf("%02d", len(req.FileUploads)) + fmt.Sprintf("%02d", req.FeeRate)
	prefix = strings.ToUpper(prefix)
	orderId := uniqueid.GenSn(prefix)

	dataSize := 0
	totalFee := int64(0)
	inscriptionRequests := make([]ordinals.InscriptionData, 0)
	for _, v := range req.FileUploads {
		// eg: dataURL, err := dataurl.DecodeString(`data:text/plain;charset=utf-8;base64,aGV5YQ==`)
		dataURL, err := dataurl.DecodeString(v.DataUrl)
		if err != nil {
			logx.Errorf("parse dataUrl error error: %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "parse dataUrl error: %v", err.Error())
		}

		if len(dataURL.Data) > 2_000_000 {
			logx.Errorf("data size too large")
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "data too large error: %v", err.Error())
		}

		// insert inscribe data into db
		_, err = l.svcCtx.TbInscribeDataModel.Insert(l.ctx, &model.TbInscribeData{
			OrderId:     orderId,
			Data:        v.DataUrl,
			ContentType: dataURL.MediaType.String(),
			FileName:    v.FileName,
		})
		if err != nil {
			logx.Errorf("insert inscribedata error %v", err.Error())
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert inscribedata error %v", err.Error())
		}

		inscriptionRequests = append(inscriptionRequests, ordinals.InscriptionData{
			ContentType: dataURL.MediaType.String(),
			Body:        dataURL.Data,
			Destination: req.ReceiveAddress,
		})

	}
	totalFee, _, err = ordinals.EstimateFee(l.svcCtx.ChainCfg, req.FeeRate, true, inscriptionRequests, int64(utxoSat))
	if err != nil {
		logx.Errorf("EstimateFee error %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "EstimateFee error %v", err.Error())

	}

	logx.Infof("==========totalFee : %v", totalFee)
	logx.Infof("==========net name: %v", l.svcCtx.ChainCfg.Name)

	// bitfish price
	if bitfishTotalPrice > 1000 {
		totalFee += bitfishTotalPrice
	}

	createTime := time.Now()
	ord := model.TbInscribeOrder{
		OrderId:        orderId,
		Count:          int64(count),
		DepositAddress: depositAddress,
		ReceiveAddress: req.ReceiveAddress,
		DataBytes:      int64(dataSize),
		FeeRate:        int64(req.FeeRate),
		TxfeeAmountSat: totalFee,
		ServiceFeeSat:  0,
		TotalAmountSat: totalFee,
		OrderStatus:    "NOTPAID",
		Version:        int64(versionNumber), // magic number for bitcoinfish order
		CreateTime:     createTime,
	}
	_, err = l.svcCtx.TbInscribeOrderModel.Insert(l.ctx, &ord)
	if err != nil {
		// TODO: use transaction ?
		l.svcCtx.TbAddressModel.Delete(l.ctx, addressId)

		logx.Errorf("insert error:%v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "insert error: %v", err.Error())
	}

	// bitcoinfish
	for _, v := range req.FileUploads {
		prefix := "bitcoinfish_"
		suffix := ".html"
		if strings.HasPrefix(v.FileName, prefix) && strings.HasSuffix(v.FileName, suffix) {
			// parse path
			b64Path, _ := strings.CutPrefix(v.FileName, prefix)
			b64Path, _ = strings.CutSuffix(b64Path, suffix)
			mergePath, _ := base64.StdEncoding.DecodeString(b64Path)

			// insert path
			l.svcCtx.TbBitfishMergePathModel.Insert(l.ctx, &model.TbBitfishMergePath{
				MergePath: string(mergePath),
			})
		}
	}

	// update bloomFilter for deposit
	if err = l.svcCtx.DepositBloomFilter.Add([]byte(depositAddress)); err != nil {
		time.Sleep(time.Millisecond * 17)
		l.svcCtx.DepositBloomFilter.Add([]byte(depositAddress))
	}

	filenames := make([]string, 0)
	for _, v := range req.FileUploads {
		filenames = append(filenames, v.FileName)
	}

	resp = &types.CreateInscribeOrderResp{
		OrderId:        ord.OrderId,
		Count:          count,
		Filenames:      filenames,
		DepositAddress: ord.DepositAddress,
		ReceiveAddress: ord.ReceiveAddress,
		FeeRate:        int(ord.FeeRate),
		Bytes:          dataSize,
		InscribeFee:    int(ord.TxfeeAmountSat),
		ServiceFee:     int(ord.ServiceFeeSat),
		Total:          int(ord.TotalAmountSat),
		CreateTime:     createTime.Unix(),
		// CreateTime:     createTime.Format("2006-01-02 15:04:05 +0800 CST"),
	}

	return
}

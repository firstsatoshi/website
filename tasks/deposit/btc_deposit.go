package deposit

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/mempool"
	"github.com/firstsatoshi/website/common/task"
	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// btc deposit task implemention

var _ task.Task = &BtcDepositTask{}

type BtcDepositTask struct {
	ctx  context.Context
	stop context.CancelFunc

	chainCfg *chaincfg.Params
	apiHost  string

	redis  *redis.Redis
	config *config.Config

	apiClient *mempool.MempoolApiClient

	tbDepositModel   model.TbDepositModel
	tbBlockscanModel model.TbBlockscanModel
	tbAddressModel   model.TbAddressModel
}

func NewBtcDepositTask(apiHost string, config *config.Config, chainCfg *chaincfg.Params) *BtcDepositTask {
	ctx, cancel := context.WithCancel(context.Background())

	redis, err := redis.NewRedis(config.CacheRedis[0].RedisConf)
	if err != nil {
		panic(err)
	}

	sqlConn := sqlx.NewMysql(config.MySql.DataSource)

	apiClient := mempool.NewMempoolApiClient(apiHost)

	return &BtcDepositTask{
		ctx:  ctx,
		stop: cancel,

		config:   config,
		redis:    redis,
		apiHost:  apiHost,
		chainCfg: chainCfg,

		apiClient: apiClient,

		tbDepositModel:   model.NewTbDepositModel(sqlConn, config.CacheRedis),
		tbBlockscanModel: model.NewTbBlockscanModel(sqlConn, config.CacheRedis),
		tbAddressModel:   model.NewTbAddressModel(sqlConn, config.CacheRedis),
	}
}

func (t *BtcDepositTask) Start() {

	for {
		ticker := time.NewTicker(time.Second * 7)
		select {
		case <-t.ctx.Done():
			logx.Info("Gracefully exit Inscribe Task goroutine....")
			// wait sub-goroutine
			return
		case <-ticker.C:
			logx.Info("======= Btc Inscribe Task =================")
			t.scanBlock()
		}
	}
}

func (t *BtcDepositTask) Stop() {
	t.stop()
}

func (t *BtcDepositTask) scanBlock() {

	// load all listen address into redis from db
	counter := 0
	maxId, err := t.tbAddressModel.FindMaxId(t.ctx, "BTC")
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return
	}
	for i := int32(0); i < maxId; i++ {
		a, err := t.tbAddressModel.FindOne(t.ctx, int64(i))
		if err != nil {
			continue
		}

		// it'll load address into redis
		_, err = t.tbAddressModel.FindOneByAddress(t.ctx, a.Address)
		if err != nil {
			continue
		}

		counter += 1
	}
	logx.Infof(" ===== load  %v address into redis", counter)

	// get latest height
	latestBlockHeight, err := t.apiClient.GetTipBlockHeight()
	if err != nil {
		logx.Errorf("GetTipBlockHeight error: %v", err.Error())
		return
	}

	// get blockHeight from db
	blockScan, err := t.tbBlockscanModel.FindOneByCoinType(t.ctx, "BTC")
	if err != nil {
		// if blockHieght doesn't exists , insert the latest height
		if err == model.ErrNotFound {
			_, err := t.tbBlockscanModel.Insert(t.ctx, &model.TbBlockscan{
				CoinType:    "BTC",
				BlockNumber: int64(latestBlockHeight),
			})
			if err != nil {
				logx.Errorf("Insert error: %v", err.Error())
				return
			}

			logx.Infof("initial insert blockscan block height %v", latestBlockHeight)
			return
		}
		logx.Errorf("FindOneByCoinType error: %v", err.Error())
		return
	}

	// if blockHeight >= latest height
	if blockScan.BlockNumber >= int64(latestBlockHeight) {
		logx.Infof("=== reached latest block %v,  sleep some time", latestBlockHeight)
		return
	}

	for n := uint64(blockScan.BlockNumber); n <= uint64(latestBlockHeight); n++ {
		// 	get block hash by height
		blockHash, err := t.apiClient.GetBlockHashByHeight(n)
		if err != nil {
			logx.Errorf("GetBlockHashByHeight error: %v", err.Error())
			return
		}

		// 	get block transactions txids by blockhash
		txids, err := t.apiClient.GetBlockTansactionIDs(blockHash)
		if err != nil {
			logx.Errorf("GetBlockTansactionIDs error: %v", err.Error())
			return
		}

		for _, txid := range txids {
			//  get each transaction details by txid
			tx, err := t.apiClient.GetTansaction(txid)
			if err != nil {
				logx.Errorf("GetTansaction error: %v", err.Error())
				return
			}

			// because we only fecth transaction from block, so this is impossible be false
			if tx.Status.Confirmed != true {
				continue
			}

			// DO NOT support locktime transaction deposit
			if tx.Locktime != 0 {
				continue
			}

			//  check vout of transaction whether it is deposit transaction by check receive address whether in redis

			// support multi order deposit in one transaction
			totalDepositValueMap := make(map[string]uint64, 0)

			for _, vo := range tx.Vout {
				if vo.ScriptpubkeyType != "v1_p2tr" {
					continue
				}

				if vo.Value <= 1000 {
					continue
				}

				addr, err := t.tbAddressModel.FindOneByAddress(t.ctx, vo.ScriptpubkeyAddress)
				if err == model.ErrNotFound {
					continue
				}
				if err != nil {
					logx.Errorf("FindOneByAddress error: %v", err.Error())
					return
				}

				if addr.Address != vo.ScriptpubkeyAddress {
					logx.Errorf("address NOT MATCH ???")
					return
				}

				// is deposit, accumulate the value
				totalDepositValueMap[addr.Address] += vo.Value
			}

			// insert into db
			for depositAddr, value := range totalDepositValueMap {
				_, err := t.tbDepositModel.Insert(t.ctx, &model.TbDeposit{
					CoinType:    "BTC",
					FromAddress: "----NO-USE---",
					ToAddress:   depositAddr,
					Txid:        txid,
					Amount:      int64(value),
					Decimals:    8,
				})

				if err != nil {
					logx.Errorf("Insert error:%v", err.Error())
					return
				}
			}

			//  if tx is deposit transaction check the value whether equal or greater than order's payment value
			// for depositAddr, value := range totalDepositValueMap {
			// }
		}
	}

}

package deposit

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/bmfilter"
	"github.com/firstsatoshi/website/common/globalvar"
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

	sqlConn sqlx.SqlConn

	bloomFilter *bmfilter.BloomFilter

	apiClient *mempool.MempoolApiClient

	tbDepositModel           model.TbDepositModel
	tbBlockscanModel         model.TbBlockscanModel
	tbAddressModel           model.TbAddressModel
	tbOrderModel             model.TbOrderModel
	tbBlindboxModel          model.TbBlindboxModel
	tbLockOrderBlindboxModel model.TbLockOrderBlindboxModel
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

		sqlConn:     sqlConn,
		apiClient:   apiClient,
		bloomFilter: bmfilter.NewUpgwBloomFilter(redis, globalvar.BTC),

		tbDepositModel:           model.NewTbDepositModel(sqlConn, config.CacheRedis),
		tbBlockscanModel:         model.NewTbBlockscanModel(sqlConn, config.CacheRedis),
		tbAddressModel:           model.NewTbAddressModel(sqlConn, config.CacheRedis),
		tbBlindboxModel:          model.NewTbBlindboxModel(sqlConn, config.CacheRedis),
		tbOrderModel:             model.NewTbOrderModel(sqlConn, config.CacheRedis),
		tbLockOrderBlindboxModel: model.NewTbLockOrderBlindboxModel(sqlConn, config.CacheRedis),
	}
}

// Start implement task.Task.Start() interface()
func (t *BtcDepositTask) Start() {

	for {
		ticker := time.NewTicker(time.Second * 7)
		select {
		case <-t.ctx.Done():
			logx.Info("Gracefully exit Inscribe Task goroutine....")
			// wait sub-goroutine
			return
		case <-ticker.C:
			logx.Info("======= Btc Deposit Task =================")
			t.scanBlock()

			//
			time.Sleep(2)
			t.txMempool()
		}
	}
}

// Stop implement task.Task.Stop() interface()
func (t *BtcDepositTask) Stop() {
	t.stop()
}

// scanBlock to poll each block and parse deposit transaction
func (t *BtcDepositTask) scanBlock() {
	// load all listen address into redis bloomfilter
	counter := 0
	addresses, err := t.tbAddressModel.FindAll(t.ctx, globalvar.BTC)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return
	}

	for i := 0; i < len(addresses); i++ {
		if err := t.bloomFilter.Add([]byte(addresses[i].Address)); err != nil {
			logx.Errorf("error: %v", err.Error())
			return
		}

		counter += 1
	}
	logx.Infof(" ===== load  %v address into redis bloom filter", counter)

	// get latest height
	latestBlockHeight, err := t.apiClient.GetTipBlockHeight()
	if err != nil {
		logx.Errorf("GetTipBlockHeight error: %v", err.Error())
		return
	}

	// get blockHeight from db
	blockScan, err := t.tbBlockscanModel.FindOneByCoinType(t.ctx, globalvar.BTC)
	if err != nil {
		// if blockHieght doesn't exists , insert the latest height
		if err == model.ErrNotFound {
			_, err := t.tbBlockscanModel.Insert(t.ctx, &model.TbBlockscan{
				CoinType:    globalvar.BTC,
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

	for n := uint64(blockScan.BlockNumber + 1); n <= uint64(latestBlockHeight); n++ {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

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

		// multi goroutine txfecther
		goroutineCount := 30
		ch := make(chan mempool.Transaction, 50000)
		go func() {
			wg := sync.WaitGroup{}
			l := len(txids)
			e := len(txids) / goroutineCount
			for i := 0; i < goroutineCount; i++ {

				startIdx := i * e
				endIdx := startIdx + e
				if i == goroutineCount-1 {
					endIdx = l
				}

				wg.Add(1)
				go t.txFecther(i, &wg, txids[startIdx:endIdx], ch)
			}
			wg.Wait()

			// close channel
			close(ch)
		}()

		okTxCount := 0
		for tx := range ch {

			okTxCount += 1

			txid := tx.Txid
			// logx.Infof("txid: %v", txid)

			// because we only fecth transaction from block, so this is impossible be false
			if tx.Status.Confirmed != true {
				logx.Infof("===skip===%v", txid)
				continue
			}

			// https://learnmeabitcoin.com/technical/locktime
			// DO NOT support locktime transaction deposit
			if tx.Locktime >= 500000000 || tx.Locktime > tx.Status.BlockHeight+1 {
				logx.Infof("locktime is %v ", tx.Locktime)
				continue
			}

			//  check vout of transaction whether it is deposit transaction by check receive address whether in redis

			// support multi order deposit in one transaction
			totalDepositValueMap := make(map[string]uint64, 0)

			for _, vo := range tx.Vout {
				if vo.ScriptpubkeyType != "v1_p2tr" {
					// logx.Infof("is not v1_p2tr address")
					continue
				}

				if vo.Value <= 1000 {
					logx.Infof("too small value")
					continue
				}

				// check address whther is bloom filter
				isExists, err := t.bloomFilter.Exists([]byte(vo.ScriptpubkeyAddress))
				if err != nil {
					logx.Errorf("bloom filter check error: %v ", err.Error())
					return
				}
				if !isExists {
					// logx.Infof("%v is not deposit address", vo.ScriptpubkeyAddress)
					continue
				}

				// check again by querying database
				addr, err := t.tbAddressModel.FindOneByAddress(t.ctx, vo.ScriptpubkeyAddress)
				if err == model.ErrNotFound {
					logx.Infof("not found %v ", vo.ScriptpubkeyAddress)
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
				_, err := t.tbDepositModel.FindOneByToAddressTxid(t.ctx, depositAddr, txid)
				if err == model.ErrNotFound {
					_, err := t.tbDepositModel.Insert(t.ctx, &model.TbDeposit{
						CoinType:    globalvar.BTC,
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
					logx.Infof("==== insert db ok! === ")
				}
			}

			// lock blindbox
			//  if tx is deposit transaction check the value whether equal or greater than order's payment value
			for depositAddr, value := range totalDepositValueMap {
				order, err := t.tbOrderModel.FindOneByDepositAddress(t.ctx, depositAddr)
				if err != nil {
					if err == model.ErrNotFound {
						logx.Errorf("=========== DEPOSIT ADDRESS NOT MATCH ORDER: %v =======", depositAddr)
						continue
					}
					logx.Errorf("FindOneByDepositAddress error: %v", err.Error())
					continue
				}

				if value < uint64(order.TotalAmountSat) {
					// TODO ??
					logx.Infof(" ============ DEPOSIT AMOUNT IS NOT ENOUGH, order total is %v, got %v ====", order.TotalAmountSat, value)
					continue
				}

				// update order
				order.PayTxid = sql.NullString{Valid: true, String: txid}
				order.PayTime = sql.NullTime{Valid: true, Time: time.Now()}
				order.OrderStatus = "PAYSUCCESS"
				order.PayConfirmedTime = sql.NullTime{Valid: true, Time: time.Now()}
				order.Version += 1
				if err := t.tbOrderModel.Update(t.ctx, order); err != nil {
					logx.Errorf("Update: %v", err.Error())
					return
				}

				// get not lock blindbox
				query := t.tbBlindboxModel.RowBuilder().Where(squirrel.Eq{
					"is_active": 1,
					"is_locked": 0,
					"status":    "NOTMINT",
				})
				boxs, err := t.tbBlindboxModel.FindBlindbox(t.ctx, query)
				if err != nil {
					logx.Errorf("FindBlindbox error: %v", err.Error())
					return
				}

				// check boxs count
				if len(boxs) < int(order.Count) {
					logx.Errorf("======== BLINDBOX NOT ENOUGH TO BE LOCKED ============")
					return
				}

				// random lock boxs
				boxIds := make([]int64, 0)
				for _, b := range boxs {
					boxIds = append(boxIds, b.Id)
				}
				rand.Seed(int64(time.Now().UnixNano()))
				rand.Shuffle(len(boxIds), func(i, j int) { boxIds[i], boxIds[j] = boxIds[j], boxIds[i] })
				lockBoxIds := boxIds[:order.Count]

				// use Transaction to lock order and blindbox
				// if could lock, lock it
				err = t.sqlConn.TransactCtx(t.ctx, func(ctx context.Context, s sqlx.Session) error {

					// lock blindbox
					for _, boxId := range lockBoxIds {
						updateBlindbox := fmt.Sprintf("UPDATE tb_blindbox SET is_locked=1 WHERE id=%v", boxId)
						result, err := s.ExecCtx(ctx, updateBlindbox)
						if err != nil {
							return err
						}
						if _, err = result.RowsAffected(); err != nil {
							return err
						}
					}

					// insert
					for _, id := range lockBoxIds {
						insertSql := fmt.Sprintf("INSERT INTO tb_lock_order_blindbox (event_id, order_id, blindbox_id) VALUES(%v, '%v', '%v')",
							order.EventId, order.OrderId, id)
						result, err := s.ExecCtx(ctx, insertSql)
						if err != nil {
							return err
						}
						if _, err = result.RowsAffected(); err != nil {
							return err
						}
					}

					return nil
				})

				if err != nil {
					logx.Errorf(" lock order and blindbox error:%v ", err.Error())
					return
				}

				logx.Infof("=== lock order and blindbox success, orderId:%v, lockBoxIds: %v ===", order.OrderId, lockBoxIds)

			}
		}

		// update block height
		if okTxCount == len(txids) {
			blockScan.BlockNumber = int64(n)
			t.tbBlockscanModel.Update(t.ctx, blockScan)
		} else {
			logx.Errorf("=======================okTxCount not equal len(txids) ==================")
			return
		}
	}

}

// txFecther to fecth transaction details
func (t *BtcDepositTask) txFecther(goroutineId int, wg *sync.WaitGroup, txids []string, ch chan<- mempool.Transaction) {
	defer wg.Done()

	for i := 0; i < len(txids); {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		txid := txids[i]
		tx, err := t.apiClient.GetTansaction(txid)
		if err != nil {
			logx.Errorf("GetTansaction error: %v", err.Error())
			continue
		}
		logx.Infof("[goroutine %v] get  txid idx : %v, DONE", goroutineId, i)

		// send to channel
		ch <- tx

		i += 1
	}
}

// txMempool to monitor bitcion mempool transaction and update order status
func (t *BtcDepositTask) txMempool() {
	// load all of 0~30 minutes order , order by create time asc
	now := time.Now()
	query := t.tbOrderModel.RowBuilder().Where(squirrel.Eq{
		"order_status": "NOTPAID",
	}).Where(squirrel.Gt{
		"create_time": time.Unix(now.Unix()-30*60, 0),
	}).Where(squirrel.Lt{
		"create_time": time.Unix(now.Unix()-1*60, 0),
	}).Limit(500).OrderBy("id DESC")

	orders, err := t.tbOrderModel.FindOrders(t.ctx, query)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return
	}
	if len(orders) == 0 {
		logx.Infof("==no order need to monitor txmempool==")
		return
	}

	// get all of address utxo by listunspent
	for _, order := range orders {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		// memTxs, err := t.apiClient.GetAddressMempoolTxs( order.DepositAddress)
		utxos, err := t.apiClient.GetAddressUTXOs(order.DepositAddress)
		if err != nil {
			logx.Errorf("GetAddressUTXOs error: %v", err.Error())
			return
		}

		if len(utxos) == 0 {
			continue
		}

		for _, utxo := range utxos {
			// if tx has be confirmed, skip it. this should be processed by scanBlock,
			if utxo.TxStatus.Confirmed {
				continue
			}

			// NOTE: we only support single utxo deposit, do not support multiple utxo deposit
			//
			// check total amount of all of utxos
			// if utxo's amount  greater or equal than order's total amount
			if utxo.Value >= uint64(order.TotalAmountSat) {
				// update order's status to PAYPENDING
				order.PayTxid = sql.NullString{Valid: true, String: utxo.Txid}
				order.PayTime = sql.NullTime{Valid: true, Time: time.Now()}
				order.OrderStatus = "PAYPENDING"

				// ignore error
				t.tbOrderModel.Update(t.ctx, order)
			}
		}
	}

	// NOTE: DO NOT lock any blindboxs until the deposit transaction be confirmed(into block)
}

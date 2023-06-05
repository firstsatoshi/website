package inscribe

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/keymanager"
	"github.com/firstsatoshi/website/common/mempool"
	"github.com/firstsatoshi/website/common/ordinals"
	"github.com/firstsatoshi/website/common/task"
	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ordinals inscribe implement

var _ task.Task = &BtcInscribeTask{}

type BtcInscribeTask struct {
	ctx  context.Context
	stop context.CancelFunc

	chainCfg *chaincfg.Params
	apiHost  string

	redis   *redis.Redis
	config  *config.Config
	sqlConn sqlx.SqlConn

	apiClient *mempool.MempoolApiClient

	changeAddressKeyManager *keymanager.KeyManager
	depositAddressKm        *keymanager.KeyManager

	tbDepositModel             model.TbDepositModel
	tbAddressModel             model.TbAddressModel
	tbOrderModel               model.TbOrderModel
	tbBlindboxModel            model.TbBlindboxModel
	tbTbLockOrderBlindboxModel model.TbLockOrderBlindboxModel
}

func NewBtcInscribeTask(apiHost string, config *config.Config, chainCfg *chaincfg.Params) *BtcInscribeTask {
	ctx, cancel := context.WithCancel(context.Background())

	redis, err := redis.NewRedis(config.CacheRedis[0].RedisConf)
	if err != nil {
		panic(err)
	}

	sqlConn := sqlx.NewMysql(config.MySql.DataSource)

	apiClient := mempool.NewMempoolApiClient(apiHost)

	if len(os.Getenv("CHANGE_SEED")) == 0 {
		panic("empty CHANGE_SEED")
	}
	changeKm, err := keymanager.NewKeyManagerFromSeed(os.Getenv("CHANGE_SEED"), *chainCfg)
	if err != nil {
		panic(err)
	}

	if len(os.Getenv("DEPOSIT_SEED")) == 0 {
		panic("empty DEPOSIT_SEED")
	}
	depositAddressKm, err := keymanager.NewKeyManagerFromSeed(os.Getenv("DEPOSIT_SEED"), *chainCfg)
	if err != nil {
		panic(err)
	}

	return &BtcInscribeTask{
		ctx:  ctx,
		stop: cancel,

		config:  config,
		redis:   redis,
		sqlConn: sqlConn,

		apiClient: apiClient,

		changeAddressKeyManager: changeKm,
		depositAddressKm:        depositAddressKm,

		apiHost:                    apiHost,
		chainCfg:                   chainCfg,
		tbDepositModel:             model.NewTbDepositModel(sqlConn, config.CacheRedis),
		tbAddressModel:             model.NewTbAddressModel(sqlConn, config.CacheRedis),
		tbOrderModel:               model.NewTbOrderModel(sqlConn, config.CacheRedis),
		tbBlindboxModel:            model.NewTbBlindboxModel(sqlConn, config.CacheRedis),
		tbTbLockOrderBlindboxModel: model.NewTbLockOrderBlindboxModel(sqlConn, config.CacheRedis),
	}
}

func (t *BtcInscribeTask) Start() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ticker := time.NewTicker(time.Second * 20)
			select {
			case <-t.ctx.Done():
				logx.Info("Gracefully exit tx monitor Task goroutine....")
				// wait sub-goroutine
				return
			case <-ticker.C:
				logx.Info("======= Btc txmonitor task======")
				t.txMonitor()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ticker := time.NewTicker(time.Second * 20)
			select {
			case <-t.ctx.Done():
				logx.Info("Gracefully exit ordertimeout Task goroutine....")
				// wait sub-goroutine
				return
			case <-ticker.C:
				logx.Info("======= Btc ordertimeout task======")
				t.orderTimeout()
			}
		}
	}()
	defer wg.Wait()

	for {
		ticker := time.NewTicker(time.Second * 7)
		select {
		case <-t.ctx.Done():
			logx.Info("Gracefully exit Inscribe Task goroutine....")
			// wait sub-goroutine
			return
		case <-ticker.C:
			logx.Info("======= Btc Inscribe Task =================")
			t.inscribe()
		}
	}

}

func (t *BtcInscribeTask) Stop() {
	t.stop()
}

func (t *BtcInscribeTask) inscribe() {

	// get order from db, 1 order per time
	query := t.tbOrderModel.RowBuilder().Where(squirrel.Eq{
		"order_status": "PAYSUCCESS",
	}).Limit(1)
	orders, err := t.tbOrderModel.FindOrders(t.ctx, query)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return
	}
	if len(orders) == 0 {
		logx.Infof("==no order need to inscribe==")
		return
	}

	order := orders[0]

	// get locked images by order
	q := t.tbTbLockOrderBlindboxModel.RowBuilder().Where(squirrel.Eq{
		"order_id": order.OrderId,
	})

	blindboxs, err := t.tbTbLockOrderBlindboxModel.FindAll(t.ctx, q, "")
	if err != nil {
		logx.Errorf("FindAll error: %v", err.Error())
		return
	}

	// get utxo
	depositAddress, err := btcutil.DecodeAddress(order.DepositAddress, t.chainCfg)
	if err != nil {
		logx.Errorf("DecodeAddress error: %v", err.Error())
		return
	}

	utxos, err := t.apiClient.ListUnspent(depositAddress)
	if err != nil {
		logx.Errorf("ListUnspent error: %v", err.Error())
		return
	}

	// check balance(utxo) - price > fee
	balanceSat := int64(0)
	for _, utxo := range utxos {
		if utxo.Output.Value > 10000 {
			balanceSat += utxo.Output.Value
		}
	}
	totalPriceSat := order.PriceSat * order.Count
	feeSat := balanceSat - totalPriceSat

	// make inscribe data
	inscribeData := make([]ordinals.InscriptionData, 0)
	for _, bbox := range blindboxs {
		imgFilePath := fmt.Sprintf("/images/%v.png", bbox.BlindboxId)
		imgData, err := ioutil.ReadFile(imgFilePath)
		if err != nil {
			logx.Errorf("ReadFile read image %v error: %v", imgFilePath, err.Error())
			return
		}

		insData := ordinals.InscriptionData{
			ContentType: "image/png",
			Body:        imgData,
			Destination: order.ReceiveAddress,
		}
		inscribeData = append(inscribeData, insData)
	}

	_, changeAddress, err := t.changeAddressKeyManager.GetWifKeyAndAddresss(0, uint32(order.Id%5))
	logx.Infof("orderId:%v, CHANGEADDRESS: %v", order.OrderId, changeAddress)

	addr, err := t.tbAddressModel.FindOneByAddress(t.ctx, order.DepositAddress)
	if err != nil {
		logx.Errorf("FindOneByAddress error:%v", order.DepositAddress)
		return
	}

	depositWif, depositAddressStr, err := t.depositAddressKm.GetWifKeyAndAddresss(uint32(addr.AccountIndex), uint32(addr.AddressIndex))
	if err != nil {
		logx.Errorf("GetWifKeyAndAddresss error: %v", err.Error())
		return
	}
	defer func() { depositWif = "" }()

	if addr.Address != depositAddressStr {
		logx.Errorf("====== DEPOSITADDRESS ADDRESS NOT MATCH %v not match %v ==========", addr.Address, depositAddressStr)
		return
	}

	onlyEstimate := true
	_, _, feeEstimate, changeSat, err := ordinals.Inscribe(changeAddress, depositWif, t.chainCfg, int(order.FeeRate), inscribeData, onlyEstimate)
	if err != nil {
		logx.Errorf(" estimate fee error: %v ", err.Error())
		return
	}

	if feeEstimate > feeSat {
		// TODO: we must estimate accuracy FEE before create order
		logx.Infof("=============== feeEstimate greater than feeSat ================")
	}

	if changeSat < order.PriceSat {
		// TODO: we must estimate accuracy FEE before create order
		logx.Infof("=============== changeSat less than PriceSat ================")
	}

	// inscrbe images
	onlyEstimate = false // push tx to blockchain
	commitTxid, revealTxids, realFee, realChange, err := ordinals.Inscribe(changeAddress, depositWif, t.chainCfg, int(order.FeeRate), inscribeData, onlyEstimate)
	if err != nil {
		logx.Errorf("inscribe error: %v ", err.Error())
		return
	}
	depositWif = ""
	logx.Infof("======= OrderId: %v inscribe finished", order.OrderId)

	// TODO:
	if len(revealTxids) != len(blindboxs) {
		logx.Errorf(" revealTxids size %v NOT MATCH blindboxs size %v ", len(revealTxids), len(blindboxs))
		return // ?????
	}

	// update order status
	// update blindbox status
	for nTry := 0; ; nTry++ {
		err = t.sqlConn.TransactCtx(t.ctx, func(ctx context.Context, s sqlx.Session) error {

			// update blindbox status to MINTING
			for i, b := range blindboxs {
				revealTxid := revealTxids[i]
				updateBlindbox := fmt.Sprintf(
					"UPDATE tb_blindbox SET status='%v',commit_txid='%v',reveal_txid='%v',real_fee_sat=%v,real_change_sat=%v WHERE id=%v",
					"MINTING", commitTxid, revealTxid, realFee, realChange, b.Id)
				result, err := s.ExecCtx(ctx, updateBlindbox)
				if err != nil {
					return err
				}
				if _, err = result.RowsAffected(); err != nil {
					return err
				}
			}

			// update order status
			if true {
				updateSql := fmt.Sprintf("UPDATE tb_order SET order_status='%v' WHERE id=%v", "INSCRIBING", order.Id)
				result, err := s.ExecCtx(t.ctx, updateSql)
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
			if nTry < 3 {
				time.Sleep(time.Duration(nTry) * time.Second)
				logx.Errorf("update order status and blindbox status error, try it later")
				continue
			}

			logx.Errorf("update order status and blindbox status error :%v ", err.Error())
		}
		break
	}
	logx.Errorf("update order %v status and blindbox status  SUCCESS ", order.OrderId)
}

func (t *BtcInscribeTask) txMonitor() {
	// TODO: if deposit tx in a orphan block ?

	// get peding txs' txids from databases
	querySql := t.tbBlindboxModel.RowBuilder().Where(squirrel.Eq{
		"status": "MINTING",
	})
	mintingBoxs, err := t.tbBlindboxModel.FindBlindbox(t.ctx, querySql)
	if err != nil {
		logx.Errorf("FindBlindbox error: %v", err.Error())
		return
	}

	// monitor tx status
	for _, mbx := range mintingBoxs {
		tx, err := t.apiClient.GetTansaction(mbx.RevealTxid.String)
		if err != nil {
			logx.Errorf("GetTansaction error: %v, continue", err.Error())
			continue
		}

		// still pending
		if !tx.Status.Confirmed {
			logx.Infof(" blindbox: %v, revealTxid:%v , still pending", mbx.Id, mbx.RevealTxid.String)
			continue
		}

		// if reveal is comfirmed
		for nTry := 0; ; nTry++ {
			err = t.sqlConn.TransactCtx(t.ctx, func(ctx context.Context, s sqlx.Session) error {

				// update blindbox status to MINT
				if true {
					updateBlindbox := fmt.Sprintf("UPDATE tb_blindbox SET status='%v' WHERE id=%v", "MINT", mbx.Id)
					result, err := s.ExecCtx(ctx, updateBlindbox)
					if err != nil {
						return err
					}
					if _, err = result.RowsAffected(); err != nil {
						return err
					}
				}

				// update order status
				if true {
					updateSql := fmt.Sprintf("UPDATE tb_order SET order_status='%v' WHERE id=%v", "ALLSUCCESS", mbx.Id)
					result, err := s.ExecCtx(t.ctx, updateSql)
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
				if nTry < 3 {
					time.Sleep(time.Duration(nTry) * time.Second)
					logx.Errorf("update order status and blindbox status error, try it later")
					continue
				}

				logx.Errorf("update order status and blindbox status error :%v ", err.Error())
			}
			break
		}
	}

	// update order status when tx be succeed
}

func (t *BtcInscribeTask) orderTimeout() {

	// 120 minutes to timeout
	now := time.Now()
	timeout := time.Unix(now.Unix()-120*60*60, 0)
	queryBuilder := t.tbOrderModel.RowBuilder().Where(squirrel.Eq{
		"order_status": "NOTPAID",
	}).Where(squirrel.Lt{
		"create_time": timeout,
	}).Limit(100)

	orders, err := t.tbOrderModel.FindOrders(t.ctx, queryBuilder)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return
	}

	for _, order := range orders {
		if now.Sub(order.CreateTime).Seconds() < 120*60*60 {
			continue
		}

		// timeout
		order.OrderStatus = "PAYTIMEOUT"
		err := t.tbOrderModel.Update(t.ctx, order)
		if err != nil {
			logx.Errorf("Update error: %v", err.Error())
			continue
		}

	}

}

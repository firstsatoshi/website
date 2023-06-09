package inscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/globalvar"
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

	depositAddressKm *keymanager.KeyManager

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

		depositAddressKm: depositAddressKm,

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
			ticker := time.NewTicker(time.Second * 60)
			select {
			case <-t.ctx.Done():
				logx.Info("Gracefully exit btcCoinPrice Task goroutine....")
				// wait sub-goroutine
				return
			case <-ticker.C:
				logx.Info("======= Btc btcCoinPrice task======")
				t.btcCoinPrice()
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
	logx.Infof("orderId: %v, orderInfo: %v", order.OrderId, order)

	// get locked images by order
	q := t.tbTbLockOrderBlindboxModel.RowBuilder().Where(squirrel.Eq{
		"order_id": order.OrderId,
	})

	lockOrderBoxs, err := t.tbTbLockOrderBlindboxModel.FindAll(t.ctx, q, "")
	if err != nil {
		logx.Errorf("FindAll error: %v", err.Error())
		return
	}

	logx.Infof("lockOrderBoxs: %v", lockOrderBoxs)

	// make inscribe data
	inscribeData := make([]ordinals.InscriptionData, 0)
	for _, bbox := range lockOrderBoxs {
		imgFilePath := fmt.Sprintf("/images/%v.png", bbox.BlindboxId)
		imgData, err := ioutil.ReadFile(imgFilePath)
		if err != nil {
			logx.Errorf("ReadFile read image %v error: %v", imgFilePath, err.Error())
			return
		}
		logx.Infof("img size: %v", len(imgData))

		insData := ordinals.InscriptionData{
			ContentType: "image/png",
			Body:        imgData,
			Destination: order.ReceiveAddress,
		}
		inscribeData = append(inscribeData, insData)
	}

	// get change address
	rndIdx := int(order.Id) % len(globalvar.MainChangeAdddress)
	changeAddress := globalvar.MainChangeAdddress[rndIdx]
	if t.chainCfg.Net == wire.TestNet3 {
		changeAddress = globalvar.TestnetChangeAddress[0]
	}

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

	// inscrbe images
	onlyEstimate := false // push tx to blockchain
	commitTxid, revealTxids, realFee, realChange, err := ordinals.Inscribe(changeAddress, depositWif, t.chainCfg, int(order.FeeRate), inscribeData, onlyEstimate)
	if err != nil {
		logx.Errorf("inscribe error: %v ", err.Error())
		return
	}
	depositWif = ""
	logx.Infof("======= OrderId: %v inscribe finished", order.OrderId)

	// TODO:
	if len(revealTxids) != len(lockOrderBoxs) {
		logx.Errorf(" revealTxids size %v NOT MATCH blindboxs size %v ", len(revealTxids), len(lockOrderBoxs))
		return // ?????
	}

	// update order status
	// update blindbox status
	for nTry := 0; ; nTry++ {
		err = t.sqlConn.TransactCtx(t.ctx, func(ctx context.Context, s sqlx.Session) error {

			// update blindbox order_status to MINTING
			for i, b := range lockOrderBoxs {
				revealTxid := revealTxids[i]
				updateBlindbox := fmt.Sprintf(
					"UPDATE tb_blindbox SET status='%v',commit_txid='%v',reveal_txid='%v' WHERE id=%v",
					"MINTING", commitTxid, revealTxid, b.BlindboxId)
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
				updateSql := fmt.Sprintf("UPDATE tb_order SET order_status='%v',real_fee_sat=%v,real_change_sat=%v WHERE id=%v",
					"MINTING", realFee, realChange, order.Id)
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
	logx.Infof("update order %v status and blindbox status  SUCCESS ", order.OrderId)
}

// txMonitor monitor tx and update order and blindbox status
func (t *BtcInscribeTask) txMonitor() {

	queryOrdersSql := t.tbOrderModel.RowBuilder().Where(squirrel.Eq{
		"order_status": "MINTING",
	})
	mintingOrders, err := t.tbOrderModel.FindOrders(t.ctx, queryOrdersSql)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return
	}

	// 1---n
	orderRevealTxMap := make(map[string]([]int), 0)
	for _, mo := range mintingOrders {
		orderRevealTxMap[mo.OrderId] = make([]int, 0)

		q := t.tbTbLockOrderBlindboxModel.RowBuilder().Where(squirrel.Eq{
			"order_id": mo.OrderId,
		})
		lobs, err := t.tbTbLockOrderBlindboxModel.FindAll(t.ctx, q, "")
		if err != nil {
			logx.Errorf("FindAll error: %v", err.Error())
			return
		}

		// push  blindbox id
		for _, lob := range lobs {
			orderRevealTxMap[mo.OrderId] = append(orderRevealTxMap[mo.OrderId], int(lob.BlindboxId))
		}
	}

	successOrderMap := make(map[string]bool, 0)
	for orderId, boxIds := range orderRevealTxMap {

		okCount := 0
		for _, boxId := range boxIds {
			mbx, err := t.tbBlindboxModel.FindOne(t.ctx, int64(boxId))
			if err != nil {
				logx.Errorf("FindOne error: %v", err.Error())
				return
			}

			// monitor tx status
			tx, err := t.apiClient.GetTansaction(mbx.RevealTxid.String)
			if err != nil {
				logx.Errorf("GetTansaction error: %v, continue", err.Error())
				continue
			}

			// TODO: if deposit tx in a orphan block ?  waiting more blocks?
			// still pending
			if !tx.Status.Confirmed {
				logx.Infof(" blindbox: %v, revealTxid:%v , still pending", mbx.Id, mbx.RevealTxid.String)
				continue
			}

			if mbx.Status != "MINT" {
				mbx.Status = "MINT"
				err = t.tbBlindboxModel.Update(t.ctx, mbx)
				if err != nil {
					logx.Errorf("update order status and blindbox status error :%v ", err.Error())
					return
				}
			}
			okCount += 1
		}

		// all boxs of order has been success, set order's status to ALLSUCCESS
		if okCount == len(boxIds) {
			successOrderMap[orderId] = true
		}
	}

	// update order's status to ALLSUCCESS
	for _, mo := range mintingOrders {
		if orderSuccess, ok := successOrderMap[mo.OrderId]; ok && orderSuccess {
			mo.OrderStatus = "ALLSUCCESS"
			t.tbOrderModel.Update(t.ctx, mo)
		}
	}

}

func (t *BtcInscribeTask) orderTimeout() {

	// 120 minutes to timeout
	timeoutSecs := 120 * 60
	now := time.Now()
	timeout := time.Unix(now.Unix()-int64(timeoutSecs), 0)
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
		if now.Sub(order.CreateTime).Seconds() < float64(timeoutSecs) {
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

func (t *BtcInscribeTask) btcCoinPrice() {
	key := "bitcion-price-usd"

	for i := 0; i < 5; i++ {

		type DataItem struct {
			PriceUsd string `json:"priceUsd"`
		}
		type Resp struct {
			Data DataItem `json:"data"`
		}

		rsp, err := http.Get("https://api.coincap.io/v2/assets/bitcoin")
		if err != nil {
			logx.Errorf("error: %v", err.Error())
			time.Sleep(time.Second)
			continue
		}

		if rsp.StatusCode != http.StatusOK {
			logx.Errorf("statusCode: %v", rsp.StatusCode)
			time.Sleep(time.Second)
			continue
		}

		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		var r Resp
		if err = json.Unmarshal(body, &r); err != nil {
			time.Sleep(time.Second)
			logx.Errorf("json Unmarshal error: %v", err.Error())
			continue
		}

		err = t.redis.Set(key, r.Data.PriceUsd)
		if err != nil {
			logx.Errorf("t.redis.Set error: %v", err.Error())
			continue
		}
		break
	}

}

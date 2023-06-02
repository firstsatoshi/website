package inscribe

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

	return &BtcInscribeTask{
		ctx:  ctx,
		stop: cancel,

		config:  config,
		redis:   redis,
		sqlConn: sqlConn,

		apiClient: apiClient,

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
	}).Limit(1)

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

	_, _, feeEstimate, changeSat, err := ordinals.Inscribe("TODO", "TODO", t.chainCfg, int(order.FeeRate), inscribeData, true)
	if err != nil {
		logx.Errorf(" estimate fee error: %v ", err.Error())
		return
	}

	if feeEstimate > feeSat {
		// TODO: we must estimate accuracy FEE before create order
	}

	if changeSat < order.PriceSat {
		// TODO: we must estimate accuracy FEE before create order
	}

	// inscrbe images
	commitTxid, revealTxids, realFee, realChange, err := ordinals.Inscribe("TODO", "TODO", t.chainCfg, int(order.FeeRate), inscribeData, true)
	if err != nil {
		logx.Errorf("estimate fee error: %v ", err.Error())
		return
	}
	logx.Infof("======= OrderId: %v inscribe finished", order.OrderId)

	// TODO:
	if len(revealTxids) != len(blindboxs) {
		logx.Errorf(" revealTxids size NOT MATCH blindboxs size ")
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
				time.Sleep(1)
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

	// get txids from databases

	// async monitor tx status

	// update order status when tx be succeed
}

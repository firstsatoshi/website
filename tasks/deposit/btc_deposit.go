package deposit

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
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

	return &BtcDepositTask{
		ctx:  ctx,
		stop: cancel,

		config:           config,
		redis:            redis,
		apiHost:          apiHost,
		chainCfg:         chainCfg,
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
		}
	}
}

func (t *BtcDepositTask) Stop() {
	t.stop()
}

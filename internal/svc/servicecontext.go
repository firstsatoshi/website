package svc

import (
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/bmfilter"
	"github.com/firstsatoshi/website/common/keymanager"
	"github.com/firstsatoshi/website/common/unisat"
	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/model"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	TbWaitlistModel          model.TbWaitlistModel
	TbBlindboxEventModel     model.TbBlindboxEventModel
	TbOrderModel             model.TbOrderModel
	TbAddressModel           model.TbAddressModel
	TbBlindboxModel          model.TbBlindboxModel
	TbLockOrderBlindboxModel model.TbLockOrderBlindboxModel
	TbInscribeOrderModel     model.TbInscribeOrderModel
	TbInscribeDataModel      model.TbInscribeDataModel

	// redis
	Redis *redis.Redis

	// deposit address bloom filter
	DepositBloomFilter *bmfilter.BloomFilter

	KeyManager *keymanager.KeyManager
	ChainCfg   *chaincfg.Params

	PeriodLimit *limit.PeriodLimit

	// unisat api client
	UnisatApiClient *unisat.UnisatApiClient
}

func NewServiceContext(c config.Config) *ServiceContext {

	// no cache, only database
	sqlConn := sqlx.NewMysql(c.MySql.DataSource)

	seed := os.Getenv("DEPOSIT_SEED")
	if len(seed) == 0 {
		panic("empty DEPOSIT_SEED")
	}

	chainCfg := chaincfg.MainNetParams
	if len(os.Getenv("BITEAGLE_TESTNET")) != 0 {
		chainCfg = chaincfg.TestNet3Params
	}
	km, err := keymanager.NewKeyManagerFromSeed(seed, chainCfg)
	if err != nil {
		panic(err)
	}

	rds := c.CacheRedis[0].NewRedis()

	periodLimit := limit.NewPeriodLimit(15, 3, rds, "api-createorder-rate-limit-key")

	return &ServiceContext{
		Config:                   c,
		Redis:                    rds,
		DepositBloomFilter:       bmfilter.NewUpgwBloomFilter(rds, "BTC"),
		TbWaitlistModel:          model.NewTbWaitlistModel(sqlConn, c.CacheRedis),
		TbBlindboxEventModel:     model.NewTbBlindboxEventModel(sqlConn, c.CacheRedis),
		TbOrderModel:             model.NewTbOrderModel(sqlConn, c.CacheRedis),
		TbAddressModel:           model.NewTbAddressModel(sqlConn, c.CacheRedis),
		TbBlindboxModel:          model.NewTbBlindboxModel(sqlConn, c.CacheRedis),
		TbLockOrderBlindboxModel: model.NewTbLockOrderBlindboxModel(sqlConn, c.CacheRedis),
		TbInscribeOrderModel:     model.NewTbInscribeOrderModel(sqlConn, c.CacheRedis),
		TbInscribeDataModel:      model.NewTbInscribeDataModel(sqlConn, c.CacheRedis),
		KeyManager:               km,
		ChainCfg:                 &chainCfg,
		PeriodLimit:              periodLimit,
		UnisatApiClient:          unisat.NewUnisatApiClient(),
	}
}

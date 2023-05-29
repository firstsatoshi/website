package svc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/fantopia-dev/website/common/keymanager"
	"github.com/fantopia-dev/website/internal/config"
	"github.com/fantopia-dev/website/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	TbWaitlistModel      model.TbWaitlistModel
	TbBlindboxEventModel model.TbBlindboxEventModel
	TbOrderModel         model.TbOrderModel
	TbAddressModel       model.TbAddressModel

	// redis
	Redis *redis.Redis

	KeyManager *keymanager.KeyManager
}

func NewServiceContext(c config.Config) *ServiceContext {

	// no cache, only database
	sqlConn := sqlx.NewMysql(c.MySql.DataSource)

	seed := "[20230529byyoungqqcn@163.com]:__.+-&2$fz&lGp)93-_-x$.-x_4.-~`_T_92fn^lsYTpz-N-"
	km, err := keymanager.NewKeyManagerFromSeed(seed, chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:               c,
		Redis:                c.CacheRedis[0].NewRedis(),
		TbWaitlistModel:      model.NewTbWaitlistModel(sqlConn, c.CacheRedis),
		TbBlindboxEventModel: model.NewTbBlindboxEventModel(sqlConn, c.CacheRedis),
		TbOrderModel:         model.NewTbOrderModel(sqlConn, c.CacheRedis),
		TbAddressModel:       model.NewTbAddressModel(sqlConn, c.CacheRedis),
		KeyManager:           km,
	}
}

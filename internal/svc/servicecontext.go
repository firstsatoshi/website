package svc

import (
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

	// redis
	Redis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	// no cache, only database
	sqlConn := sqlx.NewMysql(c.MySql.DataSource)

	return &ServiceContext{
		Config:               c,
		Redis:                c.CacheRedis[0].NewRedis(),
		TbWaitlistModel:      model.NewTbWaitlistModel(sqlConn, c.CacheRedis),
		TbBlindboxEventModel: model.NewTbBlindboxEventModel(sqlConn, c.CacheRedis),
		TbOrderModel:         model.NewTbOrderModel(sqlConn, c.CacheRedis),
	}
}

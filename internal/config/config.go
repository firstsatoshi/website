package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 数据库
	MySql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	LogConf logx.LogConf
}

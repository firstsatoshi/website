package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbWaitlistModel = (*customTbWaitlistModel)(nil)

type (
	// TbWaitlistModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbWaitlistModel.
	TbWaitlistModel interface {
		tbWaitlistModel
	}

	customTbWaitlistModel struct {
		*defaultTbWaitlistModel
	}
)

// NewTbWaitlistModel returns a model for the database table.
func NewTbWaitlistModel(conn sqlx.SqlConn, c cache.CacheConf) TbWaitlistModel {
	return &customTbWaitlistModel{
		defaultTbWaitlistModel: newTbWaitlistModel(conn, c),
	}
}

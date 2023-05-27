package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbOrderModel = (*customTbOrderModel)(nil)

type (
	// TbOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbOrderModel.
	TbOrderModel interface {
		tbOrderModel
	}

	customTbOrderModel struct {
		*defaultTbOrderModel
	}
)

// NewTbOrderModel returns a model for the database table.
func NewTbOrderModel(conn sqlx.SqlConn, c cache.CacheConf) TbOrderModel {
	return &customTbOrderModel{
		defaultTbOrderModel: newTbOrderModel(conn, c),
	}
}

package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbOrderBlindboxModel = (*customTbOrderBlindboxModel)(nil)

type (
	// TbOrderBlindboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbOrderBlindboxModel.
	TbOrderBlindboxModel interface {
		tbOrderBlindboxModel
	}

	customTbOrderBlindboxModel struct {
		*defaultTbOrderBlindboxModel
	}
)

// NewTbOrderBlindboxModel returns a model for the database table.
func NewTbOrderBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) TbOrderBlindboxModel {
	return &customTbOrderBlindboxModel{
		defaultTbOrderBlindboxModel: newTbOrderBlindboxModel(conn, c),
	}
}

package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbDepositModel = (*customTbDepositModel)(nil)

type (
	// TbDepositModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbDepositModel.
	TbDepositModel interface {
		tbDepositModel
	}

	customTbDepositModel struct {
		*defaultTbDepositModel
	}
)

// NewTbDepositModel returns a model for the database table.
func NewTbDepositModel(conn sqlx.SqlConn, c cache.CacheConf) TbDepositModel {
	return &customTbDepositModel{
		defaultTbDepositModel: newTbDepositModel(conn, c),
	}
}

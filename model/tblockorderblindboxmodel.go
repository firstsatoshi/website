package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbLockOrderBlindboxModel = (*customTbLockOrderBlindboxModel)(nil)

type (
	// TbLockOrderBlindboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbLockOrderBlindboxModel.
	TbLockOrderBlindboxModel interface {
		tbLockOrderBlindboxModel
	}

	customTbLockOrderBlindboxModel struct {
		*defaultTbLockOrderBlindboxModel
	}
)

// NewTbLockOrderBlindboxModel returns a model for the database table.
func NewTbLockOrderBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) TbLockOrderBlindboxModel {
	return &customTbLockOrderBlindboxModel{
		defaultTbLockOrderBlindboxModel: newTbLockOrderBlindboxModel(conn, c),
	}
}

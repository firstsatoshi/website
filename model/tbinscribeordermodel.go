package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbInscribeOrderModel = (*customTbInscribeOrderModel)(nil)

type (
	// TbInscribeOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbInscribeOrderModel.
	TbInscribeOrderModel interface {
		tbInscribeOrderModel
	}

	customTbInscribeOrderModel struct {
		*defaultTbInscribeOrderModel
	}
)

// NewTbInscribeOrderModel returns a model for the database table.
func NewTbInscribeOrderModel(conn sqlx.SqlConn, c cache.CacheConf) TbInscribeOrderModel {
	return &customTbInscribeOrderModel{
		defaultTbInscribeOrderModel: newTbInscribeOrderModel(conn, c),
	}
}

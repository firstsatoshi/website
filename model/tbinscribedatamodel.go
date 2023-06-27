package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbInscribeDataModel = (*customTbInscribeDataModel)(nil)

type (
	// TbInscribeDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbInscribeDataModel.
	TbInscribeDataModel interface {
		tbInscribeDataModel
	}

	customTbInscribeDataModel struct {
		*defaultTbInscribeDataModel
	}
)

// NewTbInscribeDataModel returns a model for the database table.
func NewTbInscribeDataModel(conn sqlx.SqlConn, c cache.CacheConf) TbInscribeDataModel {
	return &customTbInscribeDataModel{
		defaultTbInscribeDataModel: newTbInscribeDataModel(conn, c),
	}
}

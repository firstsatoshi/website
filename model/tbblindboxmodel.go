package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBlindboxModel = (*customTbBlindboxModel)(nil)

type (
	// TbBlindboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBlindboxModel.
	TbBlindboxModel interface {
		tbBlindboxModel
	}

	customTbBlindboxModel struct {
		*defaultTbBlindboxModel
	}
)

// NewTbBlindboxModel returns a model for the database table.
func NewTbBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) TbBlindboxModel {
	return &customTbBlindboxModel{
		defaultTbBlindboxModel: newTbBlindboxModel(conn, c),
	}
}

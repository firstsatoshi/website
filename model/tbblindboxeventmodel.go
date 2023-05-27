package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBlindboxEventModel = (*customTbBlindboxEventModel)(nil)

type (
	// TbBlindboxEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBlindboxEventModel.
	TbBlindboxEventModel interface {
		tbBlindboxEventModel
	}

	customTbBlindboxEventModel struct {
		*defaultTbBlindboxEventModel
	}
)

// NewTbBlindboxEventModel returns a model for the database table.
func NewTbBlindboxEventModel(conn sqlx.SqlConn, c cache.CacheConf) TbBlindboxEventModel {
	return &customTbBlindboxEventModel{
		defaultTbBlindboxEventModel: newTbBlindboxEventModel(conn, c),
	}
}

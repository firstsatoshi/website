package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBlockscanModel = (*customTbBlockscanModel)(nil)

type (
	// TbBlockscanModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBlockscanModel.
	TbBlockscanModel interface {
		tbBlockscanModel
	}

	customTbBlockscanModel struct {
		*defaultTbBlockscanModel
	}
)

// NewTbBlockscanModel returns a model for the database table.
func NewTbBlockscanModel(conn sqlx.SqlConn, c cache.CacheConf) TbBlockscanModel {
	return &customTbBlockscanModel{
		defaultTbBlockscanModel: newTbBlockscanModel(conn, c),
	}
}

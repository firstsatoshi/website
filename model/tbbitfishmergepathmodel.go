package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBitfishMergePathModel = (*customTbBitfishMergePathModel)(nil)

type (
	// TbBitfishMergePathModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBitfishMergePathModel.
	TbBitfishMergePathModel interface {
		tbBitfishMergePathModel
	}

	customTbBitfishMergePathModel struct {
		*defaultTbBitfishMergePathModel
	}
)

// NewTbBitfishMergePathModel returns a model for the database table.
func NewTbBitfishMergePathModel(conn sqlx.SqlConn, c cache.CacheConf) TbBitfishMergePathModel {
	return &customTbBitfishMergePathModel{
		defaultTbBitfishMergePathModel: newTbBitfishMergePathModel(conn, c),
	}
}

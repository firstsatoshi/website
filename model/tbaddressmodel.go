package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbAddressModel = (*customTbAddressModel)(nil)

type (
	// TbAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbAddressModel.
	TbAddressModel interface {
		tbAddressModel
	}

	customTbAddressModel struct {
		*defaultTbAddressModel
	}
)

// NewTbAddressModel returns a model for the database table.
func NewTbAddressModel(conn sqlx.SqlConn, c cache.CacheConf) TbAddressModel {
	return &customTbAddressModel{
		defaultTbAddressModel: newTbAddressModel(conn, c),
	}
}

package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbOrderModel = (*customTbOrderModel)(nil)

type (
	// TbOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbOrderModel.
	TbOrderModel interface {
		tbOrderModel
		RowBuilder() squirrel.SelectBuilder
		FindOrders(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbOrder, error)
	}

	customTbOrderModel struct {
		*defaultTbOrderModel
	}
)

// NewTbOrderModel returns a model for the database table.
func NewTbOrderModel(conn sqlx.SqlConn, c cache.CacheConf) TbOrderModel {
	return &customTbOrderModel{
		defaultTbOrderModel: newTbOrderModel(conn, c),
	}
}

// export to logic use
func (m *customTbOrderModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbOrderRows).From(m.table)
}

// FindOrders query order by custom sql
func (m *customTbOrderModel) FindOrders(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbOrder, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbOrder
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

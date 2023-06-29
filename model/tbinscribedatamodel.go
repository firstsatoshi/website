package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbInscribeDataModel = (*customTbInscribeDataModel)(nil)

type (
	// TbInscribeDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbInscribeDataModel.
	TbInscribeDataModel interface {
		tbInscribeDataModel
		RowBuilder() squirrel.SelectBuilder
		FindInscribeDatas(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbInscribeData, error)
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

// export to logic use
func (m *customTbInscribeDataModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbOrderRows).From(m.table)
}

// FindOrders query order by custom sql
func (m *customTbInscribeDataModel) FindInscribeDatas(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbInscribeData, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbInscribeData
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBlindboxEventModel = (*customTbBlindboxEventModel)(nil)

type (
	// TbBlindboxEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBlindboxEventModel.
	TbBlindboxEventModel interface {
		tbBlindboxEventModel

		RowBuilder() squirrel.SelectBuilder
		FindBlindboxEvents(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbBlindboxEvent, error)
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

// export to logic use
func (m *customTbBlindboxEventModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbBlindboxEventRows).From(m.table)
}

func (m *customTbBlindboxEventModel) FindBlindboxEvents(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbBlindboxEvent, error) {

	query, values, err := rowBuilder.OrderBy("id ASC").ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbBlindboxEvent
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

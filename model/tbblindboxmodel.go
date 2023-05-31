package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbBlindboxModel = (*customTbBlindboxModel)(nil)

type (
	// TbBlindboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbBlindboxModel.
	TbBlindboxModel interface {
		tbBlindboxModel

		RowBuilder() squirrel.SelectBuilder
		FindBlindbox(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbBlindbox, error)
		// FindAll(ctx context.Context, coinType string) ([]*TbAddress, error)
		// FindMaxId(ctx context.Context, coinType string) (int32, error)
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

// export to logic use
func (m *customTbBlindboxModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbBlindboxRows).From(m.table)
}

func (m *customTbBlindboxModel) FindBlindbox(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbBlindbox, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbBlindbox
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

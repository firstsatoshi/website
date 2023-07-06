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
		CountBuilder() squirrel.SelectBuilder
		FindBlindboxByIdNoCahce(ctx context.Context, id int64) (*TbBlindbox, error)
		FindBlindbox(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbBlindbox, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*TbBlindbox, error)
		FindCount(ctx context.Context) (int64, error)
		FindCountByBuilder(ctx context.Context, builder squirrel.SelectBuilder) (int64, error)
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

func (m *customTbBlindboxModel) FindBlindboxByIdNoCahce(ctx context.Context, id int64) (*TbBlindbox, error) {

	rowBuilder := m.RowBuilder().Where("id = ?", id)
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp TbBlindbox
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTbBlindboxModel) FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*TbBlindbox, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id ASC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := rowBuilder.Where("is_active = ?", 1).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbBlindbox
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// export logic
func (m *customTbBlindboxModel) CountBuilder() squirrel.SelectBuilder {
	return squirrel.Select("COUNT(id)").From(m.table)
}

func (m *customTbBlindboxModel) FindCount(ctx context.Context) (int64, error) {

	query, values, err := m.CountBuilder().Where("is_active = ?", 1).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customTbBlindboxModel) FindCountByBuilder(ctx context.Context, builder squirrel.SelectBuilder) (int64, error) {

	query, values, err := builder.Where("is_active = ?", 1).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

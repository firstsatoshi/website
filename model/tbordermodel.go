package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
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
		FindCount(ctx context.Context, counter squirrel.SelectBuilder) (int64, error)
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error)
		CountBuilder() squirrel.SelectBuilder
		SumBuilder(field string) squirrel.SelectBuilder
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

// export logic
func (m *customTbOrderModel) CountBuilder() squirrel.SelectBuilder {
	return squirrel.Select("COUNT(id)").From(m.table)
}

func (m *customTbOrderModel) FindCount(ctx context.Context, counter squirrel.SelectBuilder) (int64, error) {

	query, values, err := counter.ToSql()
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

func (m *customTbOrderModel) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error) {

	query, values, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	logx.Infof("===%v", query)

	var resp float64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// export logic
func (m *customTbOrderModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

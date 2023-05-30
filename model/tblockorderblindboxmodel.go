package model

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbLockOrderBlindboxModel = (*customTbLockOrderBlindboxModel)(nil)

type (
	// TbLockOrderBlindboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbLockOrderBlindboxModel.
	TbLockOrderBlindboxModel interface {
		tbLockOrderBlindboxModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		RowBuilder() squirrel.SelectBuilder
		CountBuilder(field string) squirrel.SelectBuilder
		SumBuilder(field string) squirrel.SelectBuilder
		// DeleteSoft(ctx context.Context, session sqlx.Session, data *TbLockOrderBlindbox) error
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*TbLockOrderBlindbox, error)
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*TbLockOrderBlindbox, error)
	}

	customTbLockOrderBlindboxModel struct {
		*defaultTbLockOrderBlindboxModel
	}
)

// NewTbLockOrderBlindboxModel returns a model for the database table.
func NewTbLockOrderBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) TbLockOrderBlindboxModel {
	return &customTbLockOrderBlindboxModel{
		defaultTbLockOrderBlindboxModel: newTbLockOrderBlindboxModel(conn, c),
	}
}

func (m *defaultTbLockOrderBlindboxModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {

	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

// func (m *defaultTbLockOrderBlindboxModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *HomestayOrder) error {
// 	data.DelState = 1
// 	data.DeleteTime = time.Now()
// 	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
// 		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "HomestayOrderModel delete err : %+v", err)
// 	}
// 	return nil
// }

func (m *defaultTbLockOrderBlindboxModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*TbLockOrderBlindbox, error) {

	query, values, err := rowBuilder.Where("deleted = ?", 0).ToSql()
	if err != nil {
		return nil, err
	}

	var resp TbLockOrderBlindbox
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *defaultTbLockOrderBlindboxModel) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error) {

	query, values, err := sumBuilder.Where("deleted = ?", 0).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultTbLockOrderBlindboxModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

	query, values, err := countBuilder.Where("deleted = ?", 0).ToSql()
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

func (m *defaultTbLockOrderBlindboxModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*TbLockOrderBlindbox, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("deleted = ?", 0).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbLockOrderBlindbox
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// export logic
func (m *defaultTbLockOrderBlindboxModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbLockOrderBlindboxRows).From(m.table)
}

// export logic
func (m *defaultTbLockOrderBlindboxModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// export logic
func (m *defaultTbLockOrderBlindboxModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

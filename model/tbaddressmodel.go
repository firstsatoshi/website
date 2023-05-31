package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbAddressModel = (*customTbAddressModel)(nil)

type (
	// TbAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbAddressModel.
	TbAddressModel interface {
		tbAddressModel

		FindAddresses(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbAddress, error)
		FindAll(ctx context.Context, coinType string) ([]*TbAddress, error)
		FindMaxId(ctx context.Context, coinType string) (int32, error)
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

// export to logic use
func (m *customTbAddressModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(tbAddressRows).From(m.table)
}

func (m *customTbAddressModel) FindAddresses(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*TbAddress, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*TbAddress
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

// TODO: as FindOne has used cache to accelerate query, so we could use FindOne to filter addresses, instead of query
func (m *customTbAddressModel) FindAll(ctx context.Context, coinType string) ([]*TbAddress, error) {

	query, values, err := m.RowBuilder().Where("token_name=?", coinType).Where("is_active=?", 1).ToSql()
	if err != nil {
		return nil, err
	}

	// var resp []*TbAddress
	// err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	// m.QueryRowCtx(ctx, &resp, query, values...)

	var resp []*TbAddress
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		return nil, err
	}

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTbAddressModel) FindMaxId(ctx context.Context, coinType string) (int32, error) {

	type Rsp struct {
		MaxIdx sql.NullInt32 `db:"maxidx"`
	}

	var sqlRet Rsp
	// query := fmt.Sprintf("select MAX(address_index) maxidx from %s where `token_name`='%s'", m.table, coinType)
	query := fmt.Sprintf("select MAX(id) maxid from %s where `token_name`='%s'", m.table, coinType)

	err := m.QueryRowNoCache(&sqlRet, query)
	if err != nil {
		return 0, err
	}
	if sqlRet.MaxIdx.Valid {
		return sqlRet.MaxIdx.Int32, nil
	}
	return 0, err // TODO: set default account start index
}

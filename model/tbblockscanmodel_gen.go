// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tbBlockscanFieldNames          = builder.RawFieldNames(&TbBlockscan{})
	tbBlockscanRows                = strings.Join(tbBlockscanFieldNames, ",")
	tbBlockscanRowsExpectAutoSet   = strings.Join(stringx.Remove(tbBlockscanFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tbBlockscanRowsWithPlaceHolder = strings.Join(stringx.Remove(tbBlockscanFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTbBlockscanIdPrefix       = "cache:tbBlockscan:id:"
	cacheTbBlockscanCoinTypePrefix = "cache:tbBlockscan:coinType:"
)

type (
	tbBlockscanModel interface {
		Insert(ctx context.Context, data *TbBlockscan) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TbBlockscan, error)
		FindOneByCoinType(ctx context.Context, coinType string) (*TbBlockscan, error)
		Update(ctx context.Context, data *TbBlockscan) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTbBlockscanModel struct {
		sqlc.CachedConn
		table string
	}

	TbBlockscan struct {
		Id          int64     `db:"id"`           // id
		CoinType    string    `db:"coin_type"`    // 地址类型,BTC,ETH,USDT
		BlockNumber int64     `db:"block_number"` // 区块高度
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 最后更新时间
	}
)

func newTbBlockscanModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTbBlockscanModel {
	return &defaultTbBlockscanModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`tb_blockscan`",
	}
}

func (m *defaultTbBlockscanModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	tbBlockscanCoinTypeKey := fmt.Sprintf("%s%v", cacheTbBlockscanCoinTypePrefix, data.CoinType)
	tbBlockscanIdKey := fmt.Sprintf("%s%v", cacheTbBlockscanIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, tbBlockscanCoinTypeKey, tbBlockscanIdKey)
	return err
}

func (m *defaultTbBlockscanModel) FindOne(ctx context.Context, id int64) (*TbBlockscan, error) {
	tbBlockscanIdKey := fmt.Sprintf("%s%v", cacheTbBlockscanIdPrefix, id)
	var resp TbBlockscan
	err := m.QueryRowCtx(ctx, &resp, tbBlockscanIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbBlockscanRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTbBlockscanModel) FindOneByCoinType(ctx context.Context, coinType string) (*TbBlockscan, error) {
	tbBlockscanCoinTypeKey := fmt.Sprintf("%s%v", cacheTbBlockscanCoinTypePrefix, coinType)
	var resp TbBlockscan
	err := m.QueryRowIndexCtx(ctx, &resp, tbBlockscanCoinTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `coin_type` = ? limit 1", tbBlockscanRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, coinType); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTbBlockscanModel) Insert(ctx context.Context, data *TbBlockscan) (sql.Result, error) {
	tbBlockscanCoinTypeKey := fmt.Sprintf("%s%v", cacheTbBlockscanCoinTypePrefix, data.CoinType)
	tbBlockscanIdKey := fmt.Sprintf("%s%v", cacheTbBlockscanIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tbBlockscanRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.CoinType, data.BlockNumber)
	}, tbBlockscanCoinTypeKey, tbBlockscanIdKey)
	return ret, err
}

func (m *defaultTbBlockscanModel) Update(ctx context.Context, newData *TbBlockscan) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	tbBlockscanCoinTypeKey := fmt.Sprintf("%s%v", cacheTbBlockscanCoinTypePrefix, data.CoinType)
	tbBlockscanIdKey := fmt.Sprintf("%s%v", cacheTbBlockscanIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbBlockscanRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.CoinType, newData.BlockNumber, newData.Id)
	}, tbBlockscanCoinTypeKey, tbBlockscanIdKey)
	return err
}

func (m *defaultTbBlockscanModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTbBlockscanIdPrefix, primary)
}

func (m *defaultTbBlockscanModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbBlockscanRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTbBlockscanModel) tableName() string {
	return m.table
}

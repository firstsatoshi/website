// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tbOrderBlindboxFieldNames          = builder.RawFieldNames(&TbOrderBlindbox{})
	tbOrderBlindboxRows                = strings.Join(tbOrderBlindboxFieldNames, ",")
	tbOrderBlindboxRowsExpectAutoSet   = strings.Join(stringx.Remove(tbOrderBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tbOrderBlindboxRowsWithPlaceHolder = strings.Join(stringx.Remove(tbOrderBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTbOrderBlindboxIdPrefix         = "cache:tbOrderBlindbox:id:"
	cacheTbOrderBlindboxBlindboxIdPrefix = "cache:tbOrderBlindbox:blindboxId:"
	cacheTbOrderBlindboxOrderIdPrefix    = "cache:tbOrderBlindbox:orderId:"
)

type (
	tbOrderBlindboxModel interface {
		Insert(ctx context.Context, data *TbOrderBlindbox) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TbOrderBlindbox, error)
		FindOneByBlindboxId(ctx context.Context, blindboxId string) (*TbOrderBlindbox, error)
		FindOneByOrderId(ctx context.Context, orderId string) (*TbOrderBlindbox, error)
		Update(ctx context.Context, data *TbOrderBlindbox) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTbOrderBlindboxModel struct {
		sqlc.CachedConn
		table string
	}

	TbOrderBlindbox struct {
		Id         int64  `db:"id"`          // id
		OrderId    string `db:"order_id"`    // 订单号
		BlindboxId string `db:"blindbox_id"` // 盲盒id
	}
)

func newTbOrderBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTbOrderBlindboxModel {
	return &defaultTbOrderBlindboxModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`tb_order_blindbox`",
	}
}

func (m *defaultTbOrderBlindboxModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	tbOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxIdPrefix, id)
	tbOrderBlindboxOrderIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxOrderIdPrefix, data.OrderId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, tbOrderBlindboxBlindboxIdKey, tbOrderBlindboxIdKey, tbOrderBlindboxOrderIdKey)
	return err
}

func (m *defaultTbOrderBlindboxModel) FindOne(ctx context.Context, id int64) (*TbOrderBlindbox, error) {
	tbOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxIdPrefix, id)
	var resp TbOrderBlindbox
	err := m.QueryRowCtx(ctx, &resp, tbOrderBlindboxIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbOrderBlindboxRows, m.table)
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

func (m *defaultTbOrderBlindboxModel) FindOneByBlindboxId(ctx context.Context, blindboxId string) (*TbOrderBlindbox, error) {
	tbOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxBlindboxIdPrefix, blindboxId)
	var resp TbOrderBlindbox
	err := m.QueryRowIndexCtx(ctx, &resp, tbOrderBlindboxBlindboxIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `blindbox_id` = ? limit 1", tbOrderBlindboxRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, blindboxId); err != nil {
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

func (m *defaultTbOrderBlindboxModel) FindOneByOrderId(ctx context.Context, orderId string) (*TbOrderBlindbox, error) {
	tbOrderBlindboxOrderIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxOrderIdPrefix, orderId)
	var resp TbOrderBlindbox
	err := m.QueryRowIndexCtx(ctx, &resp, tbOrderBlindboxOrderIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `order_id` = ? limit 1", tbOrderBlindboxRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, orderId); err != nil {
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

func (m *defaultTbOrderBlindboxModel) Insert(ctx context.Context, data *TbOrderBlindbox) (sql.Result, error) {
	tbOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxIdPrefix, data.Id)
	tbOrderBlindboxOrderIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxOrderIdPrefix, data.OrderId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tbOrderBlindboxRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.OrderId, data.BlindboxId)
	}, tbOrderBlindboxBlindboxIdKey, tbOrderBlindboxIdKey, tbOrderBlindboxOrderIdKey)
	return ret, err
}

func (m *defaultTbOrderBlindboxModel) Update(ctx context.Context, newData *TbOrderBlindbox) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	tbOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxIdPrefix, data.Id)
	tbOrderBlindboxOrderIdKey := fmt.Sprintf("%s%v", cacheTbOrderBlindboxOrderIdPrefix, data.OrderId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbOrderBlindboxRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.OrderId, newData.BlindboxId, newData.Id)
	}, tbOrderBlindboxBlindboxIdKey, tbOrderBlindboxIdKey, tbOrderBlindboxOrderIdKey)
	return err
}

func (m *defaultTbOrderBlindboxModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTbOrderBlindboxIdPrefix, primary)
}

func (m *defaultTbOrderBlindboxModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbOrderBlindboxRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTbOrderBlindboxModel) tableName() string {
	return m.table
}

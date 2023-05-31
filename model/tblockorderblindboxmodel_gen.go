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
	tbLockOrderBlindboxFieldNames          = builder.RawFieldNames(&TbLockOrderBlindbox{})
	tbLockOrderBlindboxRows                = strings.Join(tbLockOrderBlindboxFieldNames, ",")
	tbLockOrderBlindboxRowsExpectAutoSet   = strings.Join(stringx.Remove(tbLockOrderBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tbLockOrderBlindboxRowsWithPlaceHolder = strings.Join(stringx.Remove(tbLockOrderBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTbLockOrderBlindboxIdPrefix         = "cache:tbLockOrderBlindbox:id:"
	cacheTbLockOrderBlindboxBlindboxIdPrefix = "cache:tbLockOrderBlindbox:blindboxId:"
)

type (
	tbLockOrderBlindboxModel interface {
		Insert(ctx context.Context, data *TbLockOrderBlindbox) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TbLockOrderBlindbox, error)
		FindOneByBlindboxId(ctx context.Context, blindboxId int64) (*TbLockOrderBlindbox, error)
		Update(ctx context.Context, data *TbLockOrderBlindbox) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTbLockOrderBlindboxModel struct {
		sqlc.CachedConn
		table string
	}

	TbLockOrderBlindbox struct {
		Id         int64     `db:"id"`          // id
		EventId    int64     `db:"event_id"`    // 活动id
		OrderId    string    `db:"order_id"`    // 订单号
		BlindboxId int64     `db:"blindbox_id"` // 盲盒id
		Version    int64     `db:"version"`     // 版本号
		Deleted    int64     `db:"deleted"`     // 逻辑删除
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 最后更新时间
	}
)

func newTbLockOrderBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTbLockOrderBlindboxModel {
	return &defaultTbLockOrderBlindboxModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`tb_lock_order_blindbox`",
	}
}

func (m *defaultTbLockOrderBlindboxModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	tbLockOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbLockOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, tbLockOrderBlindboxBlindboxIdKey, tbLockOrderBlindboxIdKey)
	return err
}

func (m *defaultTbLockOrderBlindboxModel) FindOne(ctx context.Context, id int64) (*TbLockOrderBlindbox, error) {
	tbLockOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxIdPrefix, id)
	var resp TbLockOrderBlindbox
	err := m.QueryRowCtx(ctx, &resp, tbLockOrderBlindboxIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbLockOrderBlindboxRows, m.table)
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

func (m *defaultTbLockOrderBlindboxModel) FindOneByBlindboxId(ctx context.Context, blindboxId int64) (*TbLockOrderBlindbox, error) {
	tbLockOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxBlindboxIdPrefix, blindboxId)
	var resp TbLockOrderBlindbox
	err := m.QueryRowIndexCtx(ctx, &resp, tbLockOrderBlindboxBlindboxIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `blindbox_id` = ? limit 1", tbLockOrderBlindboxRows, m.table)
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

func (m *defaultTbLockOrderBlindboxModel) Insert(ctx context.Context, data *TbLockOrderBlindbox) (sql.Result, error) {
	tbLockOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbLockOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tbLockOrderBlindboxRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.EventId, data.OrderId, data.BlindboxId, data.Version, data.Deleted)
	}, tbLockOrderBlindboxBlindboxIdKey, tbLockOrderBlindboxIdKey)
	return ret, err
}

func (m *defaultTbLockOrderBlindboxModel) Update(ctx context.Context, newData *TbLockOrderBlindbox) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	tbLockOrderBlindboxBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxBlindboxIdPrefix, data.BlindboxId)
	tbLockOrderBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbLockOrderBlindboxRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.EventId, newData.OrderId, newData.BlindboxId, newData.Version, newData.Deleted, newData.Id)
	}, tbLockOrderBlindboxBlindboxIdKey, tbLockOrderBlindboxIdKey)
	return err
}

func (m *defaultTbLockOrderBlindboxModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTbLockOrderBlindboxIdPrefix, primary)
}

func (m *defaultTbLockOrderBlindboxModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbLockOrderBlindboxRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTbLockOrderBlindboxModel) tableName() string {
	return m.table
}

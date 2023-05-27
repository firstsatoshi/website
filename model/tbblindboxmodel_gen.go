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
	tbBlindboxFieldNames          = builder.RawFieldNames(&TbBlindbox{})
	tbBlindboxRows                = strings.Join(tbBlindboxFieldNames, ",")
	tbBlindboxRowsExpectAutoSet   = strings.Join(stringx.Remove(tbBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tbBlindboxRowsWithPlaceHolder = strings.Join(stringx.Remove(tbBlindboxFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTbBlindboxIdPrefix = "cache:tbBlindbox:id:"
)

type (
	tbBlindboxModel interface {
		Insert(ctx context.Context, data *TbBlindbox) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TbBlindbox, error)
		Update(ctx context.Context, data *TbBlindbox) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTbBlindboxModel struct {
		sqlc.CachedConn
		table string
	}

	TbBlindbox struct {
		Id          int64          `db:"id"`           // id
		Name        string         `db:"name"`         // 名称
		Description sql.NullString `db:"description"`  // 描述
		IsActive    int64          `db:"is_active"`    // 是否激活
		IsLocked    int64          `db:"is_locked"`    // 是否锁定
		IsInscribed int64          `db:"is_inscribed"` // 是否已铭刻(铭刻交易完全上链确认)
		CreateTime  time.Time      `db:"create_time"`  // 创建时间
		UpdateTime  time.Time      `db:"update_time"`  // 最后更新时间
	}
)

func newTbBlindboxModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTbBlindboxModel {
	return &defaultTbBlindboxModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`tb_blindbox`",
	}
}

func (m *defaultTbBlindboxModel) Delete(ctx context.Context, id int64) error {
	tbBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbBlindboxIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, tbBlindboxIdKey)
	return err
}

func (m *defaultTbBlindboxModel) FindOne(ctx context.Context, id int64) (*TbBlindbox, error) {
	tbBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbBlindboxIdPrefix, id)
	var resp TbBlindbox
	err := m.QueryRowCtx(ctx, &resp, tbBlindboxIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbBlindboxRows, m.table)
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

func (m *defaultTbBlindboxModel) Insert(ctx context.Context, data *TbBlindbox) (sql.Result, error) {
	tbBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbBlindboxIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tbBlindboxRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.Description, data.IsActive, data.IsLocked, data.IsInscribed)
	}, tbBlindboxIdKey)
	return ret, err
}

func (m *defaultTbBlindboxModel) Update(ctx context.Context, data *TbBlindbox) error {
	tbBlindboxIdKey := fmt.Sprintf("%s%v", cacheTbBlindboxIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbBlindboxRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Name, data.Description, data.IsActive, data.IsLocked, data.IsInscribed, data.Id)
	}, tbBlindboxIdKey)
	return err
}

func (m *defaultTbBlindboxModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTbBlindboxIdPrefix, primary)
}

func (m *defaultTbBlindboxModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbBlindboxRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTbBlindboxModel) tableName() string {
	return m.table
}

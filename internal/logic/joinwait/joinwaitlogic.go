package joinwait

import (
	"context"
	"net/mail"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinWaitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinWaitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinWaitLogic {
	return &JoinWaitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinWaitLogic) JoinWait(req *types.JoinWaitReq) (*types.JoinWaitResp, error) {
	var resp types.JoinWaitResp
	// verify email
	logx.Infof("email is %v", req.Email)
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_EMAIL_ERROR), "invalid email")
	}

	// check exits
	if _, err := l.svcCtx.TbWaitlistModel.FindOneByEventIdEmail(l.ctx, int64(0), req.Email); err == nil {
		resp.Duplicated = true
		return &resp, nil
	}

	// if not exits
	_, err = l.svcCtx.TbWaitlistModel.Insert(l.ctx, &model.TbWaitlist{
		EventId:   0,
		BtcAddress: "fansland",
		Email:     req.Email,
		MintLimit: 0, // NOTE: no
	})
	if err != nil {
		logx.Errorf("insert database error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	// id, _ := sqlRet.LastInsertId()
	resp.Duplicated = false
	return &resp, nil
}

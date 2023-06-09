package joinwaitlist

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/globalvar"
	"github.com/firstsatoshi/website/common/uniqueid"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinWaitListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinWaitListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinWaitListLogic {
	return &JoinWaitListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinWaitListLogic) JoinWaitList(req *types.JoinWaitListReq) (*types.JoinWaitListResp, error) {
	// verify cloudflare Turnstile token
	token, err := l.svcCtx.Redis.Get(fmt.Sprintf("%v:%v", globalvar.TURNSTILE_TOKEN_PREFIX, req.Token))
	if err != nil {
		if l.svcCtx.ChainCfg.Net == wire.TestNet3 {
			logx.Infof("============testnet skip token verify==============")
		} else {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_TOKEN_ERROR), "token: %v not exists", req.Token)
		}
	}
	logx.Infof("token: %v", token)

	// rate limit
	code, err := l.svcCtx.PeriodLimit.TakeCtx(l.ctx, req.Email)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "PeriodLimit.TakeCtx error: %v", err.Error())
	}
	if code != limit.Allowed {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TOO_MANY_REQUEST_ERROR), "rate limit error: %v", req.Email)
	}

	var resp types.JoinWaitListResp
	// verify email
	logx.Infof("email is %v", req.Email)
	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_EMAIL_ERROR), "invalid email")
	}

	// reference: https://unchained.com/blog/bitcoin-address-types-compared/
	// bc1p
	// e.g: bc1p3vs4447e5w0g828adhvpekqndtkpxmr04cj99zurxlqz50v9lz2q656na6
	// encoding: Bech32m
	btcAddress := req.BtcAddress
	if len(btcAddress) != 62 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCP2TRADDRESS_ERROR), "invalid bitcoin p2tr address")
	}

	// check exits
	if one, err := l.svcCtx.TbWaitlistModel.FindOneByEmail(l.ctx, req.Email); err == nil {
		resp.Duplicated = true
		resp.ReferalCode = uniqueid.GetReferalCodeById(one.Id)
		return &resp, nil
	}

	if one, err := l.svcCtx.TbWaitlistModel.FindOneByBtcAddress(l.ctx, req.BtcAddress); err == nil {
		resp.Duplicated = true
		resp.ReferalCode = uniqueid.GetReferalCodeById(one.Id)
		return &resp, nil
	}

	referalCode := ""
	if len(req.ReferalCode) > 0 {
		referalCode = req.ReferalCode
	}

	// if not exits
	sqlRet, err := l.svcCtx.TbWaitlistModel.Insert(l.ctx, &model.TbWaitlist{
		Email:      req.Email,
		BtcAddress: req.BtcAddress,
		RefereeId:  uniqueid.GetIdByReferalCode(referalCode),
	})
	if err != nil {
		logx.Errorf("insert database error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	id, _ := sqlRet.LastInsertId()
	resp.Duplicated = false
	resp.ReferalCode = uniqueid.GetReferalCodeById(id)
	return &resp, nil
}

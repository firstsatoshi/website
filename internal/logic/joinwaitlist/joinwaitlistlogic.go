package joinwaitlist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/mail"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/turnslite"
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
	ok, err := turnslite.VeifyToken(l.ctx, req.Token, l.svcCtx.Redis)
	if !ok || err != nil {
		if l.svcCtx.ChainCfg.Net == wire.TestNet3 {
			logx.Infof("============testnet skip token verify==============")
		} else {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_TOKEN_ERROR), "token: %v not exists", req.Token)
		}
	}

	// rate limit
	s := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(s[:])
	code, err := l.svcCtx.PeriodLimit.TakeCtx(l.ctx, "joinwaitlistapiperiodlimit:"+tokenHash)
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

	// check receiveAddress
	_, err = btcutil.DecodeAddress(req.BtcAddress, l.svcCtx.ChainCfg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_BTCADDRESS_ERROR), "invalid receive address %v", req.BtcAddress)
	}

	// check exits
	if one, err := l.svcCtx.TbWaitlistModel.FindOneByEventIdEmail(l.ctx, int64(req.EventId), req.Email); err == nil {
		resp.Duplicated = true
		resp.ReferalCode = uniqueid.GetReferalCodeById(one.Id)
		return &resp, nil
	}

	if one, err := l.svcCtx.TbWaitlistModel.FindOneByEventIdBtcAddress(l.ctx, int64(req.EventId), req.BtcAddress); err == nil {
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
		EventId:    int64(req.EventId),
		Email:      req.Email,
		BtcAddress: req.BtcAddress,
		RefereeId:  uniqueid.GetIdByReferalCode(referalCode),
		MintLimit:  0, // NOTE: no
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

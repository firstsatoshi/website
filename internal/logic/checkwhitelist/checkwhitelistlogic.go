package checkwhitelist

import (
	"context"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckWhitelistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckWhitelistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckWhitelistLogic {
	return &CheckWhitelistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckWhitelistLogic) CheckWhitelist(req *types.CheckWhitelistReq) (*types.CheckWhitelistResp, error) {

	event, err := l.svcCtx.TbBlindboxEventModel.FindOne(l.ctx, int64(req.EventId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByBtcAddress error: %v", err.Error())
	}
	if event.OnlyWhitelist == 0 {
		return &types.CheckWhitelistResp{
			EventId:     req.EventId,
			IsWhitelist: true, // Default set to true
		}, nil
	}

	_, err = l.svcCtx.TbWaitlistModel.FindOneByEventIdBtcAddress(l.ctx, int64(req.EventId), req.ReceiveAddress)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.CheckWhitelistResp{
				EventId:     req.EventId,
				IsWhitelist: false,
			}, nil
		}

		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByBtcAddress error: %v", err.Error())
	}

	return &types.CheckWhitelistResp{
		EventId:     req.EventId,
		IsWhitelist: true,
	}, nil
}

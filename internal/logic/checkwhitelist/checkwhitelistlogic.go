package checkwhitelist

import (
	"context"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"

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

	// _, err := l.svcCtx.TbWaitlistModel.FindOneByBtcAddress(l.ctx, req.ReceiveAddress)
	// if err != nil {
	// 	if err == model.ErrNotFound {
	// 		return &types.CheckWhitelistResp{
	// 			IsWhitelist: false,
	// 		}, nil
	// 	}

	// 	return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByBtcAddress error: %v", err.Error())
	// }

	return &types.CheckWhitelistResp{
		IsWhitelist: true,
	}, nil
}

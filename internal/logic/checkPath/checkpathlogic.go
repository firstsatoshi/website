package checkPath

import (
	"context"
	"encoding/base64"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/model"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPathLogic {
	return &CheckPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckPathLogic) CheckPath(req *types.CheckPathReq) (*types.CheckPathResp, error) {

	// parse path
	mergePath, err := base64.StdEncoding.DecodeString(req.Path)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "invalid bitfish path: %v", mergePath)
	}

	_, err = l.svcCtx.TbBitfishMergePathModel.FindOneByMergePath(l.ctx, string(mergePath))
	if err != nil {
		if err == model.ErrNotFound {
			return &types.CheckPathResp{
				IsExists: false,
			}, nil
		}

		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindOneByMergePath error: %v", err.Error())
	}

	return &types.CheckPathResp{
		IsExists: true,
	}, nil
}

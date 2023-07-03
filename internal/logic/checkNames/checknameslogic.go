package checkNames

import (
	"context"
	"strings"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckNamesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckNamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckNamesLogic {
	return &CheckNamesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckNamesLogic) CheckNames(req *types.CheckNamesReq) (resp []types.CheckNamesResp, err error) {

	if req.Type != "sats" {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.GET_NAMES_INFO_ERROR), "get names info error: %v", err.Error())
	}

	// .sats rules : https://docs.sats.id/sats-names/sns-spec/mint-names
	names := make([]string, 0)
	for _, name := range req.Names {
		// no space
		logx.Infof("name: %v", name)
		if strings.Contains(name, " ") {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_NAME_ERROR), "invalid name error: %v", "space is not permitted in name")
		}
		// TODO: unicode , emoji

		// all lowercase
		names = append(names, strings.ToLower(name))
	}

	checks, err := l.svcCtx.UnisatApiClient.CheckNames(req.Type, names)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.INVALID_NAME_TYPE_ERROR), "invalid name type error: %v", err.Error())
	}

	resp = make([]types.CheckNamesResp, 0)
	for k, v := range checks {
		resp = append(resp, types.CheckNamesResp{
			Name:     k,
			IsExists: v,
		})
	}
	return
}

package queryGalleryList

import (
	"context"
	"math"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryGalleryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryGalleryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryGalleryListLogic {
	return &QueryGalleryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryGalleryListLogic) QueryGalleryList(req *types.QueryGalleryListReq) (resp *types.QueryGalleryListResp, err error) {

	if req.PageSize <= 0 || req.CurPage < 0 || req.Category == "" {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "invalid request params")
	}

	builder := l.svcCtx.TbBlindboxModel.CountBuilder().Where("category = ?", req.Category)
	total, err := l.svcCtx.TbBlindboxModel.FindCountByBuilder(l.ctx, builder)
	if err != nil {
		logx.Errorf("FindCount error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindCount error")
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.PageSize)))

	if req.CurPage > totalPage {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR), "invalid request params PageSize ")
	}

	rowBuilder := l.svcCtx.TbBlindboxModel.RowBuilder().Where(squirrel.Eq{
		"category": strings.ToLower(req.Category),
	})
	boxs, err := l.svcCtx.TbBlindboxModel.FindPageListByPage(l.ctx, rowBuilder, int64(req.CurPage), int64(req.PageSize), "")
	if err != nil {
		logx.Errorf("FindPageListByPage error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "FindPageListByPage error")
	}

	resp = &types.QueryGalleryListResp{}
	resp.Category = req.Category
	resp.CurPage = req.CurPage
	resp.PageSize = req.PageSize
	resp.Total = totalPage

	for _, box := range boxs {
		resp.NFTs = append(resp.NFTs, types.NFT{
			Id:          int(box.Id),
			Name:        box.Name,
			Description: box.Description,
			ImageUrl:    box.ImgUrl,
		})
	}
	return
}

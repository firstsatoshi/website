package queryBlindboxEvent

import (
	"context"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/firstsatoshi/website/internal/types"
	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBlindboxEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryBlindboxEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBlindboxEventLogic {
	return &QueryBlindboxEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryBlindboxEventLogic) QueryBlindboxEvent() (resp *types.QueryBlindboxEventResp, err error) {
	// TODO: config by database
	EVENT_ID := 1

	event, err := l.svcCtx.TbBlindboxEventModel.FindOne(l.ctx, int64(EVENT_ID))
	if err != nil {
		logx.Errorf("TbBlindboxEventModel.FindOne error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	resp = &types.QueryBlindboxEventResp{
		EventId:      int(event.Id),
		Name:         event.EventName,
		Description:  event.EventDescription,
		PriceBtcSats: int(event.BtcPrice),
		PriceUsd:     0, // TODO
		PaymentCoin:  event.PaymentCoin,
		Supply:       int(event.Supply),
		Avail:        int(event.Avail),
		Enable:       event.IsActive > 0,
		OnlyWhiteist: event.OnlyWhitelist > 0,
		StartTime:    event.StartTime.String(),
		EndTime:      event.EndTime.String(),
	}

	return
}

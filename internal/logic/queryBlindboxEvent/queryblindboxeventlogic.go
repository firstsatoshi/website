package queryBlindboxEvent

import (
	"context"

	"github.com/Masterminds/squirrel"
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

func (l *QueryBlindboxEventLogic) QueryBlindboxEvent() (resp []types.BlindboxEvent, err error) {
	querySql := l.svcCtx.TbBlindboxEventModel.RowBuilder().Where(squirrel.Eq{
		"is_active": 1,
	})
	events, err := l.svcCtx.TbBlindboxEventModel.FindBlindboxEvents(l.ctx, querySql)
	if err != nil {
		logx.Errorf("TbBlindboxEventModel.FindOne error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	resp = make([]types.BlindboxEvent, 0)
	for _, event := range events {
		resp = append(resp, types.BlindboxEvent{
			EventId:           int(event.Id),
			Name:              event.EventName,
			Description:       event.EventDescription,
			PriceBtcSats:      int(event.PriceSats),
			PriceUsd:          0, // TODO
			MintLimit:         int(event.MintLimit),
			PaymentToken:      event.PaymentToken,
			AverageImageBytes: int(event.AverageImageBytes),
			Supply:            int(event.Supply),
			Avail:             int(event.Avail),
			Enable:            event.IsActive > 0,
			OnlyWhiteist:      event.OnlyWhitelist > 0,
			StartTime:         event.StartTime.String(),
			EndTime:           event.EndTime.String(),
		})
	}

	return
}

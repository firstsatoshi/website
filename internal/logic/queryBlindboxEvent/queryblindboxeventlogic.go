package queryBlindboxEvent

import (
	"context"
	"strings"

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

func (l *QueryBlindboxEventLogic) QueryBlindboxEvent(req *types.QueryBlindboxEventReq) (resp []types.BlindboxEvent, err error) {
	if req.EventId == 0 {
		logx.Infof("query all events")
	} else {
		logx.Infof("query event: %v", req.EventId)
	}

	querySql := l.svcCtx.TbBlindboxEventModel.RowBuilder()
	if len(req.EventStatus) > 0 {
		if req.EventStatus == "active" {
			querySql = querySql.Where(squirrel.Eq{
				"is_active": 1,
			})
		} else if req.EventStatus == "inactive" {
			querySql = querySql.Where(squirrel.Eq{
				"is_active": 0,
			})
		}
	}

	if req.EventId >= 1 {
		querySql = querySql.Where(squirrel.Eq{
			"id": req.EventId,
		})
	}

	events, err := l.svcCtx.TbBlindboxEventModel.FindBlindboxEvents(l.ctx, querySql)
	if err != nil {
		logx.Errorf("TbBlindboxEventModel.FindOne error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	// query PAYPENDING order as pendingOrders, avail = event.avail - len(pendingOrders)
	queryBuilder := l.svcCtx.TbOrderModel.RowBuilder().Where(squirrel.Eq{
		"order_status": "PAYPENDING",
	}).OrderBy("id DESC")
	if req.EventId >= 1 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{
			"event_id": req.EventId,
		})
	}

	payPendingOrders, err := l.svcCtx.TbOrderModel.FindOrders(l.ctx, queryBuilder)
	if err != nil {
		logx.Errorf("FindOrders error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SERVER_COMMON_ERROR), "database error")
	}

	resp = make([]types.BlindboxEvent, 0)
	for _, event := range events {

		// calculate safety avail
		orderCounter := 0
		for _, o := range payPendingOrders {
			if o.EventId == event.Id {
				orderCounter += 1
			}
		}
		safeAvail := event.Avail - int64(orderCounter)
		if safeAvail < 0 {
			safeAvail = 0
		}

		roadmapList := strings.Split(event.RoadmapList, ";")
		imagesList := strings.Split(event.ImgUrlList, ";")

		// parse plan list
		mintPlanList := make([]types.MintPlan, 0)
		if true {
			ss := strings.Split(event.MintPlanList, ";")
			for _, s := range ss {
				x := strings.Split(s, ",")
				if len(x) == 2 {
					mintPlanList = append(mintPlanList, types.MintPlan{
						Title: x[0],
						Plan:  x[1],
					})
				}
			}
		}

		resp = append(resp, types.BlindboxEvent{
			EventId:            int(event.Id),
			Name:               event.EventName,
			Description:        event.EventDescription,
			Detail:             event.Detail,
			AvatarImageUrl:     event.AvatarImgUrl,
			BackgroundImageUrl: event.BackgroundImgUrl,
			RoadmapDescription: event.EventDescription,
			RoadmapList:        roadmapList,
			MintPlanList:       mintPlanList,
			WebsiteUrl:         event.WebsiteUrl,
			TwitterUrl:         event.TwitterUrl,
			DiscordUrl:         event.DiscordUrl,
			ImagesList:         imagesList,
			PriceBtcSats:       int(event.PriceSats),
			PriceUsd:           0, // TODO
			MintLimit:          int(event.MintLimit),
			PaymentToken:       event.PaymentToken,
			AverageImageBytes:  int(event.AverageImageBytes),
			Supply:             int(event.Supply),
			Avail:              int(safeAvail),
			Active:             event.IsActive > 0,
			Display:            event.IsDisplay > 0,
			OnlyWhiteist:       event.OnlyWhitelist > 0,
			StartTime:          event.StartTime.Unix(),
			EndTime:            event.EndTime.Unix(),
		})
	}

	return
}

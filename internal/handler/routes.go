// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	checkwhitelist "github.com/firstsatoshi/website/internal/handler/checkwhitelist"
	coinPrice "github.com/firstsatoshi/website/internal/handler/coinPrice"
	createInscribeOrder "github.com/firstsatoshi/website/internal/handler/createInscribeOrder"
	createOrder "github.com/firstsatoshi/website/internal/handler/createOrder"
	joinwaitlist "github.com/firstsatoshi/website/internal/handler/joinwaitlist"
	queryAddress "github.com/firstsatoshi/website/internal/handler/queryAddress"
	queryBlindboxEvent "github.com/firstsatoshi/website/internal/handler/queryBlindboxEvent"
	queryGalleryList "github.com/firstsatoshi/website/internal/handler/queryGalleryList"
	queryOrder "github.com/firstsatoshi/website/internal/handler/queryOrder"
	"github.com/firstsatoshi/website/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/joinwaitlist",
				Handler: joinwaitlist.JoinWaitListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/createorder",
				Handler: createOrder.CreateOrderHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/queryorder",
				Handler: queryOrder.QueryOrderHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/queryblindboxevent",
				Handler: queryBlindboxEvent.QueryBlindboxEventHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/coinprice",
				Handler: coinPrice.CoinPriceHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/querygallerylist",
				Handler: queryGalleryList.QueryGalleryListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/checkwhitelist",
				Handler: checkwhitelist.CheckWhitelistHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/queryaddress",
				Handler: queryAddress.QueryAddressHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/createinscribeorder",
				Handler: createInscribeOrder.CreateInscribeOrderHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)
}

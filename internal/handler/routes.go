// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	createOrder "github.com/fantopia-dev/website/internal/handler/createOrder"
	joinwaitlist "github.com/fantopia-dev/website/internal/handler/joinwaitlist"
	queryBlindboxEvent "github.com/fantopia-dev/website/internal/handler/queryBlindboxEvent"
	queryOrder "github.com/fantopia-dev/website/internal/handler/queryOrder"
	"github.com/fantopia-dev/website/internal/svc"

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
				Method:  http.MethodGet,
				Path:    "/queryblindboxevent",
				Handler: queryBlindboxEvent.QueryBlindboxEventHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)
}

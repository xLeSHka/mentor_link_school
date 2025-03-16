package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	"go.uber.org/fx"
	"net/http"
)

type WebSocket struct {
	Conn      *websocket.Conn
	Clients   map[uuid.UUID]*websocket.Conn
	Broadcast chan *Message
}

func New() *WebSocket {
	var clients = make(map[uuid.UUID]*websocket.Conn)
	var broadcast = make(chan *Message)

	return &WebSocket{
		Conn:      nil,
		Clients:   clients,
		Broadcast: broadcast,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Route struct {
	routers   *ApiRouters.ApiRouters
	validator *Validators.Validator
	wsconn    *WebSocket
}

type FxOpts struct {
	fx.In
	ApiRouter *ApiRouters.ApiRouters
	Validator *Validators.Validator
	Wsconn    *WebSocket
}

func WsRoute(opts FxOpts) *Route {
	router := &Route{
		routers:   opts.ApiRouter,
		validator: opts.Validator,
		wsconn:    opts.Wsconn,
	}
	opts.ApiRouter.UserRoute.GET("/ws", router.wsconn.wsHandler)
	go opts.Wsconn.Echo()
	return router
}

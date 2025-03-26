package ws

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	"go.uber.org/fx"
	"net/http"
)

type WebSocket struct {
	Conn         *websocket.Conn
	Clients      map[uuid.UUID]*websocket.Conn
	Broadcast    chan *Message
	Consumer     *broker.Consumer
	Api          *tgbotapi.BotAPI
	UsersService service.UsersService
}
type WsFxOpts struct {
	fx.In
	Consumer     *broker.Consumer
	Api          *tgbotapi.BotAPI
	UsersService service.UsersService
}

func New(opts WsFxOpts) *WebSocket {
	var clients = make(map[uuid.UUID]*websocket.Conn)
	var broadcast = make(chan *Message)

	return &WebSocket{
		Conn:         nil,
		Clients:      clients,
		Broadcast:    broadcast,
		Consumer:     opts.Consumer,
		Api:          opts.Api,
		UsersService: opts.UsersService,
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
	opts.ApiRouter.UserRoute.GET("/ws", router.wsconn.WsHandler)
	go opts.Wsconn.Echo()
	go opts.Wsconn.WriteMessage()
	return router
}

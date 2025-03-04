package publicRoute

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gorm.io/gorm"
)

type Route struct {
	Routers *ApiRouters.ApiRouters
	DB      *gorm.DB
}

func PublicRoute(apiRouters *ApiRouters.ApiRouters, db *gorm.DB, wsconn *ws.WebSocket) *Route {
	router := &Route{
		Routers: apiRouters,
		DB:      db,
	}

	apiRouters.Public.GET("/ping", router.ping)
	apiRouters.Public.GET("/pong", func(context *gin.Context) {
		//wsconn.WriteMessage(&ws.Message{
		//	Type:    "request",
		//	UserID:  uuid.MustParse("cb168f71-fd8a-4d0d-b12a-bd91053b2fcf"),
		//	Request: &ws.Request{},
		//})
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	apiRouters.Public.GET("/mock", router.mocks)
	apiRouters.Public.StaticFile("/docsstatic/doc.json", "docs/swagger.json")
	apiRouters.Public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(config *ginSwagger.Config) {
		config.URL = "/api/docsstatic/doc.json"
	}))
	apiRouters.UserPrivate.GET("/ws", wsconn.WsHandler)
	go wsconn.Echo()

	return router
}

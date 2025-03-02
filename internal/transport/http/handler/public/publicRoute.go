package publicRoute

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.prodcontest.ru/team-14/lotti/docs"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	"gorm.io/gorm"
)

type Route struct {
	Routers *ApiRouters.ApiRouters
	DB      *gorm.DB
}

func PublicRoute(apiRouters *ApiRouters.ApiRouters, db *gorm.DB) *Route {
	router := &Route{
		Routers: apiRouters,
		DB:      db,
	}

	apiRouters.Public.GET("/ping", router.ping)
	apiRouters.Public.GET("/mock", router.mocks)
	apiRouters.Public.StaticFile("/docsstatic/doc.json", "docs/swagger.json")
	apiRouters.Public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(config *ginSwagger.Config) {
		config.URL = "/api/docsstatic/doc.json"
	}))
	apiRouters.Public.GET("/ws", router.Websocket)
	return router
}

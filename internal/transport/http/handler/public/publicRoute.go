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

func PublicRoute(apiRouters *ApiRouters.ApiRouters) *Route {
	router := &Route{
		Routers: apiRouters,
	}

	apiRouters.Public.GET("/ping", router.ping)
	apiRouters.Public.GET("/mock", router.mocks)
	apiRouters.Public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

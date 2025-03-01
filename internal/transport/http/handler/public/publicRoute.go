package publicRoute

import (
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
	return router
}

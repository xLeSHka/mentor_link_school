package publicRoute

import (
	"prodapp/internal/transport/http/handler/ApiRouters"
)

type Route struct {
	Routers *ApiRouters.ApiRouters
}

func PublicRoute(apiRouters *ApiRouters.ApiRouters) *Route {
	router := &Route{
		Routers: apiRouters,
	}

	apiRouters.Public.GET("/ping", router.ping)

	return router
}

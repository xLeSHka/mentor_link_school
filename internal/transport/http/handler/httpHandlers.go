package httpHandlers

import (
	"go.uber.org/fx"
	"prodapp/internal/app/Validators"
	"prodapp/internal/transport/http/handler/ApiRouters"
	"prodapp/internal/transport/http/handler/public"
	"prodapp/internal/transport/http/handler/user"
)

var HttpHandlers = fx.Module("httpHandlers",
	fx.Provide(
		ApiRouters.CreateApiRoutes,
		Validators.New,
		fx.Private),
	//publicRoute.PublicRoute,
	fx.Invoke(
		publicRoute.PublicRoute,
		usersRoute.UsersRoute,
	),
)

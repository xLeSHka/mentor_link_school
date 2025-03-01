package httpHandlers

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	mentorsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/mentor"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/public"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/user"
	"go.uber.org/fx"
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
		mentorsRoute.MentorsRoute,
	),
)

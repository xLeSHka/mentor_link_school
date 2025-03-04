package httpHandlers

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	groupsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/group"
	mentorsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/mentor"
	publicRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/public"
	usersRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/user"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"go.uber.org/fx"
)

var HttpHandlers = fx.Module("httpHandlers",
	fx.Provide(
		ApiRouters.CreateApiRoutes,
		Validators.New,
		ws.New,
		fx.Private),
	//publicRoute.PublicRoute,
	fx.Invoke(
		publicRoute.PublicRoute,
		usersRoute.UsersRoute,
		mentorsRoute.MentorsRoute,
		groupsRoute.GroupsRoutes,
	),
)

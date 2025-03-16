package httpHandlers

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	groupsRoute "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/group"
	mentorsRoute "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/mentor"
	publicRoute "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/public"
	studentsRoute "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/student"
	usersRoute "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/user"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
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
		studentsRoute.StudentsRoute,
		mentorsRoute.MentorsRoute,
		groupsRoute.GroupsRoutes,
	),
)
var WSHandler = fx.Module("wsHandler",
	fx.Provide(
		ApiRouters.CreateApiRoutes,
		Validators.New,
		broker.NewConsumer,
		ws.New,
		fx.Private),
	fx.Invoke(
		ws.WsRoute,
	),
)

package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"

	"go.uber.org/fx"
)

type Route struct {
	routers         *ApiRouters.ApiRouters
	validator       *Validators.Validator
	groupService    service.GroupService
	usersService    service.UserService
	minioRepository repository.MinioRepository
	wsconn          *ws.WebSocket
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	GroupService    service.GroupService
	UsersService    service.UserService
	MinioRepository repository.MinioRepository
	Ws              *ws.WebSocket
}

func GroupsRoutes(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		groupService:    opts.GroupService,
		usersService:    opts.UsersService,
		minioRepository: opts.MinioRepository,
		wsconn:          opts.Ws,
	}

	opts.ApiRouter.UserPrivate.POST("/groups/create", router.createGroup)
	opts.ApiRouter.UserPrivate.POST("/groups/:id/inviteCode", router.updateInviteCode)
	opts.ApiRouter.UserPrivate.GET("/groups/:id/members", router.getMembers)
	opts.ApiRouter.UserPrivate.POST("/groups/:id/members/role", router.updateRole)
	opts.ApiRouter.UserPrivate.GET("/groups/:id/stat", router.getStat)
	opts.ApiRouter.UserPrivate.POST("/groups/join/:code", router.acceptedInvite)
	opts.ApiRouter.UserPrivate.POST("/groups/:id/uploadAvatar", router.uploadAvatar)
	opts.ApiRouter.UserPrivate.PUT("/groups/:id/edit", router.edit)
	return router
}

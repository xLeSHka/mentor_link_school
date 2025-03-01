package groupsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"

	"go.uber.org/fx"
)

type Route struct {
	routers         *ApiRouters.ApiRouters
	validator       *Validators.Validator
	groupService    service.GroupService
	minioRepository repository.MinioRepository
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	GroupService    service.GroupService
	MinioRepository repository.MinioRepository
}

func GroupsRoutes(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		groupService:    opts.GroupService,
		minioRepository: opts.MinioRepository,
	}

	opts.ApiRouter.UserPrivate.POST("/createGroup", router.createGroup)
	opts.ApiRouter.UserPrivate.GET("", router.getGroups)
	opts.ApiRouter.GroupPrivate.GET("", router.getGroup)
	return router
}

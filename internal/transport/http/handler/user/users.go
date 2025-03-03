package usersRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"

	"go.uber.org/fx"
)

type Route struct {
	routers         *ApiRouters.ApiRouters
	validator       *Validators.Validator
	usersService    service.UserService
	minioRepository repository.MinioRepository
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	UsersService    service.UserService
	MinioRepository repository.MinioRepository
}

func UsersRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		usersService:    opts.UsersService,
		minioRepository: opts.MinioRepository,
	}

	opts.ApiRouter.UserPrivate.POST("/user/requests", router.createRequest)
	opts.ApiRouter.Public.POST("/user/auth/sign-in", router.login)
	opts.ApiRouter.UserPrivate.GET("/user/availableMentors", router.availableMentors)
	opts.ApiRouter.UserPrivate.GET("/user/mentors", router.getMyMentors)

	opts.ApiRouter.UserPrivate.GET("/init", router.init)
	opts.ApiRouter.UserPrivate.GET("/user/profile/:id", router.profileOther)
	opts.ApiRouter.UserPrivate.POST("/user/profile/edit", router.edit)
	opts.ApiRouter.UserPrivate.GET("/user/requests", router.getRequests)
	opts.ApiRouter.UserPrivate.POST("/user/uploadAvatar", router.uploadAvatar)
	//opts.ApiRouter.UserPrivate.POST("/user/invite", router.acceptedInvite)

	opts.ApiRouter.UserPrivate.GET("/ws", ws.WsHandler)
	return router
}

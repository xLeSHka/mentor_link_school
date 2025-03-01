package usersRoute

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

	opts.ApiRouter.UserPrivate.POST("/group/requests", router.createRequest)
	opts.ApiRouter.Public.POST("/user/auth/sign-in", router.login)
	opts.ApiRouter.UserPrivate.GET("/user/avaliableMentors", router.availableMentors)
	opts.ApiRouter.UserPrivate.GET("/user/mentors", router.getMyMentors)

	opts.ApiRouter.UserPrivate.GET("/user/profile", router.profile)
	opts.ApiRouter.UserPrivate.GET("/user/requests", router.getRequests)
	opts.ApiRouter.UserPrivate.POST("/user/uploadAvatar", router.uploadAvatar)

	return router
}

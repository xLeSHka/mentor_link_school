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

	opts.ApiRouter.Public.POST("/user/auth/sign-up", router.signup)
	opts.ApiRouter.Public.POST("/user/auth/sign-in", router.signin)

	opts.ApiRouter.UserPrivate.GET("/profile", router.profile)
	opts.ApiRouter.UserPrivate.POST("/uploadAvatar", router.uploadAvatar)
	return router
}

package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	"go.uber.org/fx"
)

type Route struct {
	routers         *ApiRouters.ApiRouters
	validator       *Validators.Validator
	usersService    service.UserService
	minioRepository repository.MinioRepository
	producer        *broker.Producer
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	UsersService    service.UserService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
}

func UsersRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		usersService:    opts.UsersService,
		minioRepository: opts.MinioRepository,
		producer:        opts.Producer,
	}

	opts.ApiRouter.UserPrivate.POST("/user/requests", router.createRequest)
	opts.ApiRouter.Public.POST("/user/auth/sign-in", router.login)
	opts.ApiRouter.UserPrivate.GET("/user/availableMentors", router.availableMentors)
	opts.ApiRouter.UserPrivate.GET("/user/mentors", router.getMyMentors)

	opts.ApiRouter.UserPrivate.GET("/init", router.init)
	opts.ApiRouter.UserPrivate.GET("/user/profile/:id", router.profileOther)
	opts.ApiRouter.UserPrivate.PUT("/user/profile/edit", router.edit)
	opts.ApiRouter.UserPrivate.GET("/user/requests", router.getRequests)
	opts.ApiRouter.UserPrivate.POST("/user/uploadAvatar", router.uploadAvatar)
	//opts.ApiRouter.UserPrivate.POST("/user/invite", router.acceptedInvite)

	return router
}

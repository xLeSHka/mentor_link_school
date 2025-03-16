package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/app/Validators"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ApiRouters"
	"go.uber.org/fx"
)

type Route struct {
	routers         *ApiRouters.ApiRouters
	validator       *Validators.Validator
	usersService    service.UsersService
	groupService    service.GroupService
	minioRepository repository.MinioRepository
	producer        *broker.Producer
	cryptoKey       []byte
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	UsersService    service.UsersService
	GroupService    service.GroupService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
	Config          config.Config
}

func UsersRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		usersService:    opts.UsersService,
		minioRepository: opts.MinioRepository,
		groupService:    opts.GroupService,
		producer:        opts.Producer,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}

	opts.ApiRouter.Public.POST("/users/auth/login", router.login)
	opts.ApiRouter.Public.POST("/users/auth/register", router.register)
	opts.ApiRouter.UserRoute.GET("/profile", router.profile)
	opts.ApiRouter.UserRoute.PATCH("/profile/edit", router.edit)
	opts.ApiRouter.UserRoute.POST("/uploadAvatar", router.uploadAvatar)
	opts.ApiRouter.UserRoute.GET("/groups", router.getGroups)
	opts.ApiRouter.UserRoute.GET("/profile/:id", router.profileOther)
	opts.ApiRouter.UserRoute.POST("/join/:code", router.acceptedInvite)

	//opts.ApiRouter.StudentRoute.POST("/user/invite", router.acceptedInvite)

	return router
}

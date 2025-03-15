package groupsRoute

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
	groupService    service.GroupService
	usersService    service.UserService
	minioRepository repository.MinioRepository
	producer        *broker.Producer
	cryptoKey       []byte
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	GroupService    service.GroupService
	UsersService    service.UserService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
	Config          config.Config
}

func GroupsRoutes(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		groupService:    opts.GroupService,
		usersService:    opts.UsersService,
		minioRepository: opts.MinioRepository,
		producer:        opts.Producer,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}

	opts.ApiRouter.UserPrivate.POST("/create", router.createGroup)
	opts.ApiRouter.UserPrivate.POST("/:id/inviteCode", router.updateInviteCode)
	opts.ApiRouter.UserPrivate.GET("/:id/members", router.getMembers)
	opts.ApiRouter.UserPrivate.POST("/:id/members/role", router.updateRole)
	opts.ApiRouter.UserPrivate.GET("/:id/stat", router.getStat)
	opts.ApiRouter.UserPrivate.POST("/join/:code", router.acceptedInvite)
	opts.ApiRouter.UserPrivate.POST("/:id/uploadAvatar", router.uploadAvatar)
	opts.ApiRouter.UserPrivate.PATCH("/:id/edit", router.edit)
	return router
}

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
	studentsService service.StudentService
	usersService    service.UsersService
	minioRepository repository.MinioRepository
	producer        *broker.Producer
	cryptoKey       []byte
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	GroupService    service.GroupService
	StudentsService service.StudentService
	UsersService    service.UsersService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
	Config          config.Config
}

func GroupsRoutes(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		groupService:    opts.GroupService,
		studentsService: opts.StudentsService,
		minioRepository: opts.MinioRepository,
		usersService:    opts.UsersService,
		producer:        opts.Producer,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}

	opts.ApiRouter.GroupPrivate.POST("/inviteCode", router.updateInviteCode)
	opts.ApiRouter.GroupPrivate.POST("/members/{userID}/role", router.addRole)
	opts.ApiRouter.GroupPrivate.DELETE("/members/{userID}/role", router.removeRole)
	opts.ApiRouter.GroupPrivate.GET("/members/{userID}", router.getRoles)
	opts.ApiRouter.GroupPrivate.POST("/uploadAvatar", router.uploadAvatar)
	opts.ApiRouter.UserRoute.POST("/groups/create", router.createGroup)
	opts.ApiRouter.GroupPrivate.GET("/members", router.getMembers)
	opts.ApiRouter.GroupPrivate.GET("/stat", router.getStat)
	opts.ApiRouter.GroupPrivate.PATCH("/edit", router.edit)
	return router
}

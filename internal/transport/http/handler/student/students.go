package studentsRoute

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
	studentsService service.StudentService
	usersService    service.UsersService
	groupService    service.GroupService
	minioRepository repository.MinioRepository

	producer  *broker.Producer
	cryptoKey []byte
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	StudentsService service.StudentService
	GroupService    service.GroupService
	MinioRepository repository.MinioRepository
	UserService     service.UsersService
	Producer        *broker.Producer
	Config          config.Config
}

func StudentsRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		studentsService: opts.StudentsService,
		minioRepository: opts.MinioRepository,
		groupService:    opts.GroupService,
		usersService:    opts.UserService,
		producer:        opts.Producer,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}

	opts.ApiRouter.StudentRoute.GET("/availableMentors", router.availableMentors)
	opts.ApiRouter.StudentRoute.POST("/:userID/requests", router.createRequest)
	opts.ApiRouter.StudentRoute.GET("/requests", router.getRequests)
	opts.ApiRouter.StudentRoute.GET("/mentors", router.getMyMentors)
	//opts.ApiRouter.StudentRoute.POST("/user/invite", router.acceptedInvite)

	return router
}

package mentorsRoute

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
	mentorService   service.MentorService
	studentsService service.StudentService
	minioRepository repository.MinioRepository
	usersService    service.UsersService
	producer        *broker.Producer
	cryptoKey       []byte
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	StudentsService service.StudentService
	MentorService   service.MentorService
	UsersService    service.UsersService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
	Config          config.Config
}

func MentorsRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		studentsService: opts.StudentsService,
		mentorService:   opts.MentorService,
		minioRepository: opts.MinioRepository,
		usersService:    opts.UsersService,
		producer:        opts.Producer,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}

	opts.ApiRouter.MentorRoute.GET("/students", router.students)
	opts.ApiRouter.MentorRoute.GET("/requests", router.getRequests)
	opts.ApiRouter.MentorRoute.POST("/requests", router.updateRequest)

	return router
}

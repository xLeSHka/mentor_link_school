package mentorsRoute

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
	mentorService   service.MentorService
	userService     service.UserService
	minioRepository repository.MinioRepository
	producer        *broker.Producer
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	UsersService    service.UserService
	MentorService   service.MentorService
	MinioRepository repository.MinioRepository
	Producer        *broker.Producer
}

func MentorsRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		userService:     opts.UsersService,
		mentorService:   opts.MentorService,
		minioRepository: opts.MinioRepository,
		producer:        opts.Producer,
	}

	opts.ApiRouter.MentorRoute.GET("/mentors/students", router.students)
	opts.ApiRouter.MentorRoute.GET("/mentors/requests", router.getRequests)
	opts.ApiRouter.MentorRoute.POST("/mentors/requests", router.updateRequest)

	return router
}

package mentorsRoute

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
	mentorService   service.MentorService
	userService     service.UserService
	minioRepository repository.MinioRepository
}

type FxOpts struct {
	fx.In
	ApiRouter       *ApiRouters.ApiRouters
	Validator       *Validators.Validator
	UsersService    service.UserService
	MentorService   service.MentorService
	MinioRepository repository.MinioRepository
}

func MentorsRoute(opts FxOpts) *Route {
	router := &Route{
		routers:         opts.ApiRouter,
		validator:       opts.Validator,
		userService:     opts.UsersService,
		mentorService:   opts.MentorService,
		minioRepository: opts.MinioRepository,
	}

	//opts.ApiRouter.GroupPrivate.GET("/mentors", router.getMentors)
	//opts.ApiRouter.GroupPrivate.POST("/getMentor", router.getMentorRequest)
	//opts.ApiRouter.GroupPrivate.POST("/createMentor", router.createMentorRequest)
	//opts.ApiRouter.MentorRoute.GET("/", router.getMentor)
	return router
}

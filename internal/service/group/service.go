package groupService

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"

	"go.uber.org/fx"
)

type GroupsService struct {
	groupRepository repository.GroupRepository
	minioRepository repository.MinioRepository
	userRepository  repository.UsersRepository
	cacheRepository repository.CacheRepository
}

type FxOpts struct {
	fx.In
	GroupRepository repository.GroupRepository
	MinioRepository repository.MinioRepository
	UserRepository  repository.UsersRepository
	CacheRepository repository.CacheRepository
}

func New(opts FxOpts) service.GroupService {
	return &GroupsService{
		groupRepository: opts.GroupRepository,
		minioRepository: opts.MinioRepository,
		userRepository:  opts.UserRepository,
		cacheRepository: opts.CacheRepository,
	}
}

package groupService

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"

	"go.uber.org/fx"
)

type GroupsService struct {
	groupRepository repository.GroupRepository
	minioRepository repository.MinioRepository
	userRepository  repository.UsersRepository
}

type FxOpts struct {
	fx.In
	GroupRepository repository.GroupRepository
	MinioRepository repository.MinioRepository
	UserRepository  repository.UsersRepository
}

func New(opts FxOpts) service.GroupService {
	return &GroupsService{
		groupRepository: opts.GroupRepository,
		minioRepository: opts.MinioRepository,
		userRepository:  opts.UserRepository,
	}
}

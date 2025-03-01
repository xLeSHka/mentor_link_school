package groupService

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"time"

	"go.uber.org/fx"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type GroupsService struct {
	groupRepository repository.GroupRepository
	minioRepository repository.MinioRepository
}

type FxOpts struct {
	fx.In
	GroupRepository repository.GroupRepository
	MinioRepository repository.MinioRepository
}

func New(opts FxOpts) service.GroupService {
	return &GroupsService{
		groupRepository: opts.GroupRepository,
		minioRepository: opts.MinioRepository,
	}
}

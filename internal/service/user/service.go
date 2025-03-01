package userService

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

type UsersService struct {
	usersRepository repository.UsersRepository
	minioRepository repository.MinioRepository
	groupRepository repository.GroupRepository
	jwt             *jwt.JWT
	rdb             *redis.Client
	cryptoKey       []byte
}

type FxOpts struct {
	fx.In
	UsersRepository repository.UsersRepository
	JWT             *jwt.JWT
	RDB             *redis.Client
	MinioRepository repository.MinioRepository
	GroupRepository repository.GroupRepository
	Config          config.Config
}

func New(opts FxOpts) service.UserService {
	return &UsersService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		rdb:             opts.RDB,
		minioRepository: opts.MinioRepository,
		groupRepository: opts.GroupRepository,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}
}

func (s *UsersService) GenerateAccessToken(id uuid.UUID) (string, error) {
	return s.jwt.CreateToken(jwtlib.MapClaims{
		"id": id,
	}, time.Now().Add(time.Hour*6))
}

package userService

import (
	"prodapp/internal/pkg/config"
	"prodapp/internal/repository"
	"prodapp/internal/service"
	"prodapp/internal/transport/http/pkg/jwt"
	"time"

	"go.uber.org/fx"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UsersService struct {
	usersRepository repository.UsersRepository
	minioRepository repository.MinioRepository
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
	Config          config.Config
}

func NewUsersService(opts FxOpts) service.UserService {
	return &UsersService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		rdb:             opts.RDB,
		minioRepository: opts.MinioRepository,
		cryptoKey:       []byte(opts.Config.CryptoKey),
	}
}

func (s *UsersService) GenerateAccessToken(id uuid.UUID) (string, error) {
	return s.jwt.CreateToken(jwtlib.MapClaims{
		"type": "user",
		"id":   id,
	}, time.Now().Add(time.Hour*6))
}

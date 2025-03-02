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
)

type UsersService struct {
	usersRepository  repository.UsersRepository
	minioRepository  repository.MinioRepository
	mentorRepository repository.MentorRepository
	jwt              *jwt.JWT
	//rdb             *redis.Client
	cryptoKey []byte
}

type FxOpts struct {
	fx.In
	UsersRepository repository.UsersRepository
	JWT             *jwt.JWT
	//RDB             *redis.Client
	MentorRepository repository.MentorRepository
	MinioRepository  repository.MinioRepository

	Config config.Config
}

func New(opts FxOpts) service.UserService {
	return &UsersService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		//rdb:             opts.RDB,
		minioRepository:  opts.MinioRepository,
		mentorRepository: opts.MentorRepository,
		cryptoKey:        []byte(opts.Config.CryptoKey),
	}
}

func (s *UsersService) GenerateAccessToken(id uuid.UUID) (string, error) {
	return s.jwt.CreateToken(jwtlib.MapClaims{
		"id": id,
	}, time.Now().Add(time.Hour*24*7))
}

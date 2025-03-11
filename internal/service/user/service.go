package userService

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
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
	cryptoKey        []byte
}

type FxOpts struct {
	fx.In
	UsersRepository  repository.UsersRepository
	JWT              *jwt.JWT
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

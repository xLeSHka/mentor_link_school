package mentorService

import (
	"time"

	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MentorService struct {
	usersRepository  repository.UsersRepository
	minioRepository  repository.MinioRepository
	mentorRepository repository.MentorRepository
	jwt              *jwt.JWT
	cache            repository.CacheRepository
	cryptoKey        []byte
}

type FxOpts struct {
	fx.In
	UsersRepository  repository.UsersRepository
	JWT              *jwt.JWT
	MinioRepository  repository.MinioRepository
	MentorRepository repository.MentorRepository
	Cache            repository.CacheRepository
	Config           config.Config
}

func New(opts FxOpts) service.MentorService {
	return &MentorService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		//rdb:             opts.RDB,
		minioRepository: opts.MinioRepository,
		//groupRepository:  opts.GroupRepository,
		mentorRepository: opts.MentorRepository,
		cache:            opts.Cache,
		cryptoKey:        []byte(opts.Config.CryptoKey),
	}
}

func (s *MentorService) GenerateAccessToken(id uuid.UUID) (string, error) {
	return s.jwt.CreateToken(jwtlib.MapClaims{
		"id": id,
	}, time.Now().Add(time.Hour*24*7))
}

package mentorService

import (
	"time"

	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MentorService struct {
	usersRepository repository.UsersRepository
	minioRepository repository.MinioRepository
	//groupRepository  repository.GroupRepository
	mentorRepository repository.MentorRepository
	jwt              *jwt.JWT
	//rdb              *redis.Client
	cryptoKey []byte
}

type FxOpts struct {
	fx.In
	UsersRepository repository.UsersRepository
	JWT             *jwt.JWT
	//RDB             *redis.Client
	MinioRepository repository.MinioRepository
	//GroupRepository  repository.GroupRepository
	MentorRepository repository.MentorRepository
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
		cryptoKey:        []byte(opts.Config.CryptoKey),
	}
}

func (s *MentorService) GenerateAccessToken(id uuid.UUID) (string, error) {
	return s.jwt.CreateToken(jwtlib.MapClaims{
		"type": "user",
		"id":   id,
	}, time.Now().Add(time.Hour*6))
}

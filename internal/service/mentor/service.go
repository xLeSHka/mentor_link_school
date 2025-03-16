package mentorService

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"
)

type MentorService struct {
	usersRepository   repository.UsersRepository
	minioRepository   repository.MinioRepository
	mentorRepository  repository.MentorRepository
	studentRepository repository.StudentRepository
	jwt               *jwt.JWT
	cache             repository.CacheRepository
	cryptoKey         []byte
}

type FxOpts struct {
	fx.In
	UsersRepository   repository.UsersRepository
	JWT               *jwt.JWT
	MinioRepository   repository.MinioRepository
	MentorRepository  repository.MentorRepository
	Cache             repository.CacheRepository
	StudentRepository repository.StudentRepository
	Config            config.Config
}

func New(opts FxOpts) service.MentorService {
	return &MentorService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		//rdb:             opts.RDB,
		minioRepository: opts.MinioRepository,
		//groupRepository:  opts.GroupRepository,
		mentorRepository:  opts.MentorRepository,
		studentRepository: opts.StudentRepository,
		cache:             opts.Cache,
		cryptoKey:         []byte(opts.Config.CryptoKey),
	}
}

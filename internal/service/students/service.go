package studentService

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"go.uber.org/fx"
)

type StudentService struct {
	usersRepository   repository.UsersRepository
	minioRepository   repository.MinioRepository
	mentorRepository  repository.MentorRepository
	studentRepository repository.StudentRepository
	jwt               *jwt.JWT
	cryptoKey         []byte
	cache             repository.CacheRepository
}

type FxOpts struct {
	fx.In
	UsersRepository  repository.UsersRepository
	JWT              *jwt.JWT
	MentorRepository repository.MentorRepository
	MinioRepository  repository.MinioRepository

	StudentRepository repository.StudentRepository
	Cache             repository.CacheRepository
	Config            config.Config
}

func New(opts FxOpts) service.StudentService {
	return &StudentService{
		usersRepository: opts.UsersRepository,
		jwt:             opts.JWT,
		//rdb:             opts.RDB,
		studentRepository: opts.StudentRepository,
		minioRepository:   opts.MinioRepository,
		mentorRepository:  opts.MentorRepository,
		cache:             opts.Cache,
		cryptoKey:         []byte(opts.Config.CryptoKey),
	}
}

package app

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository/cache"
	repositoryGroup "github.com/xLeSHka/mentorLinkSchool/internal/repository/group"
	repositoryMentor "github.com/xLeSHka/mentorLinkSchool/internal/repository/mentor"
	repositoryMinio "github.com/xLeSHka/mentorLinkSchool/internal/repository/minio"
	repositoryStudent "github.com/xLeSHka/mentorLinkSchool/internal/repository/student"
	repositoryUser "github.com/xLeSHka/mentorLinkSchool/internal/repository/user"

	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUser.New,
	repositoryMentor.New,
	repositoryGroup.New,
	repositoryMinio.New,
	repositoryStudent.New,
	cache.New,
)

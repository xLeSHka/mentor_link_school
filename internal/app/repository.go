package app

import (
	repositoryGroup "gitlab.prodcontest.ru/team-14/lotti/internal/repository/group"
	repositoryMentor "gitlab.prodcontest.ru/team-14/lotti/internal/repository/mentor"
	repositoryMinio "gitlab.prodcontest.ru/team-14/lotti/internal/repository/minio"
	repositoryUser "gitlab.prodcontest.ru/team-14/lotti/internal/repository/user"

	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUser.New,
	repositoryMentor.New,
	repositoryGroup.New,
	repositoryMinio.New,
)

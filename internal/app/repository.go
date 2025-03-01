package app

import (
	repositoryMentor "gitlab.prodcontest.ru/team-14/lotti/internal/repository/mentor"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/minio"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/user"

	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUser.New,
	repositoryMentor.New,

	repositoryMinio.New,
)

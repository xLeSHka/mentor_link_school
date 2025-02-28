package app

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/minio"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/user"

	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUser.NewUsersRepository,
	repositoryMinio.NewMinioRepository,
)

package app

import (
	"prodapp/internal/repository/minio"
	"prodapp/internal/repository/user"

	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUser.NewUsersRepository,
	repositoryMinio.NewMinioRepository,
)

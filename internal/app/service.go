package app

import (
	userService "gitlab.prodcontest.ru/team-14/lotti/internal/service/user"

	"go.uber.org/fx"
)

var Services = fx.Provide(
	userService.NewUsersService,
)

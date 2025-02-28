package app

import (
	userService "prodapp/internal/service/user"

	"go.uber.org/fx"
)

var Services = fx.Provide(
	userService.NewUsersService,
)

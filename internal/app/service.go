package app

import (
	groupService "gitlab.prodcontest.ru/team-14/lotti/internal/service/group"
	mentorService "gitlab.prodcontest.ru/team-14/lotti/internal/service/mentor"
	userService "gitlab.prodcontest.ru/team-14/lotti/internal/service/user"

	"go.uber.org/fx"
)

var Services = fx.Provide(
	userService.New,
	mentorService.New,
	groupService.New,
)

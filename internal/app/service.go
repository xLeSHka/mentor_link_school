package app

import (
	groupService "github.com/xLeSHka/mentorLinkSchool/internal/service/group"
	mentorService "github.com/xLeSHka/mentorLinkSchool/internal/service/mentor"
	userService "github.com/xLeSHka/mentorLinkSchool/internal/service/user"

	"go.uber.org/fx"
)

var Services = fx.Provide(
	userService.New,
	mentorService.New,
	groupService.New,
)

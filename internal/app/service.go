package app

import (
	groupService "github.com/xLeSHka/mentorLinkSchool/internal/service/group"
	mentorService "github.com/xLeSHka/mentorLinkSchool/internal/service/mentor"
	studentService "github.com/xLeSHka/mentorLinkSchool/internal/service/students"
	userService "github.com/xLeSHka/mentorLinkSchool/internal/service/users"

	"go.uber.org/fx"
)

var Services = fx.Provide(
	userService.New,
	studentService.New,
	mentorService.New,
	groupService.New,
)

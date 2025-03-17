package tg

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/tg/botkit"
	"github.com/xLeSHka/mentorLinkSchool/internal/tg/botkit/conn"
	"go.uber.org/fx"
)

var TGHandlers = fx.Module("tgHandlers",
	fx.Provide(
		conn.New,
		botkit.New,
	),
	fx.Invoke(
		Run,
	),
)

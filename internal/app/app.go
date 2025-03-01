package app

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/connetions/db"
	"gitlab.prodcontest.ru/team-14/lotti/internal/connetions/minio"
	"gitlab.prodcontest.ru/team-14/lotti/internal/integrations"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http"
	httpHandlers "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		db.New,
		http.New,
		jwt.New,
		//redis.New,
		minio.New,
	),
	integrations.Integrations,
	Repositories,
	Services,
	httpHandlers.HttpHandlers,
)

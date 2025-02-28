package app

import (
	"prodapp/internal/connetions/db"
	"prodapp/internal/connetions/minio"
	"prodapp/internal/connetions/redis"
	"prodapp/internal/integrations"
	"prodapp/internal/transport/http"
	httpHandlers "prodapp/internal/transport/http/handler"
	"prodapp/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		db.New,
		http.New,
		jwt.New,
		redis.New,
		minio.New,
	),
	integrations.Integrations,
	Repositories,
	Services,
	httpHandlers.HttpHandlers,
)

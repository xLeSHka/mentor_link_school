package app

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/db"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/minio"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/rdb"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http"
	httpHandlers "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		db.New,
		http.New,
		jwt.New,
		//redis.New,
		minio.New,
		rdb.New,
		broker.NewProducer,
	),
	Repositories,
	Services,
	httpHandlers.HttpHandlers,
)
var WSApp = fx.Options(
	fx.Provide(
		jwt.New,
		http.New,
	),
	httpHandlers.WSHandler,
)

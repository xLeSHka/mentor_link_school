package ApiRouters

import (
	"github.com/gin-gonic/gin"
	middlewares "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/middleware"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

type ApiRouters struct {
	Public       *gin.RouterGroup
	MentorRoute  *gin.RouterGroup
	UserPrivate  *gin.RouterGroup
	GroupPrivate *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine, jwt *jwt.JWT) *ApiRouters {

	publicRoute := gin.Group("/api")

	groupRoute := publicRoute.Group("")
	groupRoute.Use(middlewares.Auth(jwt))

	mentorRoute := groupRoute.Group("")
	mentorRoute.Use(middlewares.Auth(jwt))

	userRoute := publicRoute.Group("")
	userRoute.Use(middlewares.Auth(jwt))
	return &ApiRouters{
		Public:       publicRoute,
		MentorRoute:  mentorRoute,
		UserPrivate:  userRoute,
		GroupPrivate: groupRoute,
	}
}

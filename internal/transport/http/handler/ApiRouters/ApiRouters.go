package ApiRouters

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/middleware"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
)

type ApiRouters struct {
	Public       *gin.RouterGroup
	MentorRoute  *gin.RouterGroup
	UserPrivate  *gin.RouterGroup
	GroupPrivate *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine, db *gorm.DB, rdb *redis.Client, jwt *jwt.JWT) *ApiRouters {

	publicRoute := gin.Group("")

	groupRoute := publicRoute.Group("/:groupId")
	groupRoute.Use(middlewares.Auth(db, rdb, jwt, "owner"))

	mentorRoute := groupRoute.Group("/:mentorId")
	mentorRoute.Use(middlewares.Auth(db, rdb, jwt, "mentor"))

	userRoute := publicRoute.Group("/user")
	userRoute.Use(middlewares.Auth(db, rdb, jwt, "user"))
	return &ApiRouters{
		Public:       publicRoute,
		MentorRoute:  mentorRoute,
		UserPrivate:  userRoute,
		GroupPrivate: groupRoute,
	}
}

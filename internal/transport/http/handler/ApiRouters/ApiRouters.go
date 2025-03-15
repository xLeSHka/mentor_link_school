package ApiRouters

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	middlewares "github.com/xLeSHka/mentorLinkSchool/internal/transport/http/middleware"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
)

type ApiRouters struct {
	Public       *gin.RouterGroup
	MentorRoute  *gin.RouterGroup
	UserPrivate  *gin.RouterGroup
	GroupPrivate *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine, jwt *jwt.JWT, rdb *redis.Client, db *gorm.DB) *ApiRouters {

	publicRoute := gin.Group("/api")

	groupRoute := publicRoute.Group("/groups")
	groupRoute.Use(middlewares.Auth(jwt, rdb, db, "owner"))

	mentorRoute := publicRoute.Group("/mentors")
	mentorRoute.Use(middlewares.Auth(jwt, rdb, db, "mentor"))

	userRoute := publicRoute.Group("/users")
	userRoute.Use(middlewares.Auth(jwt, rdb, db, "student"))
	return &ApiRouters{
		Public:       publicRoute,
		MentorRoute:  mentorRoute,
		UserPrivate:  userRoute,
		GroupPrivate: groupRoute,
	}
}

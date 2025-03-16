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
	StudentRoute *gin.RouterGroup
	GroupPrivate *gin.RouterGroup
	UserRoute    *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine, jwt *jwt.JWT, rdb *redis.Client, db *gorm.DB) *ApiRouters {

	publicRoute := gin.Group("/api")

	userRoute := publicRoute.Group("/users")
	userRoute.Use(middlewares.Auth(jwt, rdb))

	groupRoute := publicRoute.Group("/groups/:groupID")
	groupRoute.Use(middlewares.RoleBasedAuth(jwt, rdb, db, "owner"))

	mentorRoute := groupRoute.Group("/mentors")
	mentorRoute.Use(middlewares.RoleBasedAuth(jwt, rdb, db, "mentor"))

	studentsRoute := groupRoute.Group("/students")
	studentsRoute.Use(middlewares.RoleBasedAuth(jwt, rdb, db, "student"))
	return &ApiRouters{
		Public:       publicRoute,
		MentorRoute:  mentorRoute,
		StudentRoute: studentsRoute,
		GroupPrivate: groupRoute,
		UserRoute:    userRoute,
	}
}

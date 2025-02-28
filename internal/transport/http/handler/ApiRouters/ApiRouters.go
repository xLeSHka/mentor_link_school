package ApiRouters

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"prodapp/internal/transport/http/middleware"
	"prodapp/internal/transport/http/pkg/jwt"
)

type ApiRouters struct {
	Public         *gin.RouterGroup
	CompanyPrivate *gin.RouterGroup
	UserPrivate    *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine, rdb *redis.Client, jwt *jwt.JWT) *ApiRouters {
	gin.MaxMultipartMemory = 1 << 20
	publicRoute := gin.Group("/api")

	companyRoute := publicRoute.Group("/business")
	companyRoute.Use(middlewares.Auth(rdb, jwt, "company"))

	userRoute := publicRoute.Group("/user")
	userRoute.Use(middlewares.Auth(rdb, jwt, "user"))
	publicRoute.Static("/uploads/", "./uploads/")
	return &ApiRouters{
		Public:         publicRoute,
		CompanyPrivate: companyRoute,
		UserPrivate:    userRoute,
	}
}

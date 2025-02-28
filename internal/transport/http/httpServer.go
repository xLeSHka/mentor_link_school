package http

import (
	"context"
	"fmt"
	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the request's origin matches the allowed origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func New(config config.Config, lc fx.Lifecycle) *gin.Engine {
	webServer := gin.Default()

	webServer.Use(CORSMiddleware())

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					err := webServer.Run(fmt.Sprintf(config.ServerAddress))
					if err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: nil,
		},
	)

	return webServer
}

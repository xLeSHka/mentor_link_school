package http

import (
	"context"
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"

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

// @title  			GetMentor API
// @version         1.0
// @description     GetMentor API docs

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host prod-team-14-mkg8u20m.final.prodcontest.ru
// @BasePath /

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

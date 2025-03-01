package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
	"strings"
)

func Auth(jwt *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
			c.Abort()
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			httpError.New(http.StatusUnauthorized, "Header bad format").SendError(c)
			c.Abort()
			return
		}
		data, err := jwt.VerifyToken(splitToken[1])
		if err != nil {
			httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
			c.Abort()
			return
		}

		personId := data["id"].(string)

		c.Set("personId", personId)
		c.Next()
	}
}

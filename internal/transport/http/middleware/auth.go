package middlewares

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Auth(rdb *redis.Client, jwt *jwt.JWT) gin.HandlerFunc {
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
		iat := int64(data["iat"].(float64))

		val, err := rdb.Get(context.Background(), "jwt:"+personId).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
		} else {
			resetTime, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
			println(iat, resetTime)
			if time.UnixMicro(iat).Before(time.UnixMicro(resetTime)) {
				httpError.New(http.StatusUnauthorized, "Token invalid").SendError(c)
				c.Abort()
				return
			}
		}

		c.Set("personId", personId)
		c.Next()
	}
}

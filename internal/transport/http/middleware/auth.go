package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Auth(jwt *jwt.JWT, rdb *redis.Client, db *gorm.DB, role string) gin.HandlerFunc {
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
			httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
			c.Abort()
			return
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
			is, err := rdb.SIsMember(context.Background(), "roles:"+personId, role).Result()
			if err != nil || !is {
				var r models.Role
				err = db.Model(&models.Role{}).Where("user_id = ? AND role = ?", personId, role).Find(&r).Error
				if err != nil {
					httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
					c.Abort()
					return
				}
				rdb.SAdd(context.Background(), "roles:"+personId, role)
			}
		}
		c.Set("personId", personId)
		c.Next()
	}
}

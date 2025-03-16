package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func RoleBasedAuth(jwt2 *jwt.JWT, rdb *redis.Client, db *gorm.DB, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID, err := jwt.ParseGroupID(c)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
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
		data, err := jwt2.VerifyToken(splitToken[1])
		if err != nil {
			httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
			c.Abort()
			return
		}

		personId := data["id"].(string)
		val, err := rdb.Get(context.Background(), "jwt:"+personId).Result()
		if err != nil {
			httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
			c.Abort()
			return
		} else {
			if val != splitToken[1] {
				httpError.New(http.StatusUnauthorized, "Token expired").SendError(c)
				c.Abort()
				return
			}
			is, err := rdb.SIsMember(context.Background(), "roles:"+personId+"_"+groupID.String(), role).Result()
			log.Println("Select role from redis for ", "roles:"+personId+"_"+groupID.String(), "role", role)
			if err != nil || !is {
				var r models.Role
				err = db.Model(&models.Role{}).Where("user_id = ? AND role = ? AND group_id = ?", personId, role, groupID).Find(&r).Error
				if err != nil {
					httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
					c.Abort()
					return
				}
				log.Println("Select role from postgres for ", "roles:"+personId+"_"+groupID.String(), "role", role)
				rdb.SAdd(context.Background(), "roles:"+personId+"_"+groupID.String(), role)
			}
		}
		c.Set("personId", personId)
		c.Next()
	}
}

func Auth(jwt *jwt.JWT, rdb *redis.Client) gin.HandlerFunc {
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
		val, err := rdb.Get(context.Background(), "jwt:"+personId).Result()
		if err != nil {
			httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
			c.Abort()
			return
		} else {
			if val != splitToken[1] {
				httpError.New(http.StatusUnauthorized, "Token expired").SendError(c)
				c.Abort()
				return
			}
		}
		c.Set("personId", personId)
		c.Next()
	}
}

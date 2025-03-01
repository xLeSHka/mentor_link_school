package middlewares

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type GetMentorID struct {
	ID      string `uri:"mentorId" binding:"required,uuid"`
	GroupID string `uri:"groupId" binding:"required,uuid"`
}

func Auth(db *gorm.DB, rdb *redis.Client, jwt *jwt.JWT, tokenType string) gin.HandlerFunc {
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
		if tokenType == "user" {
			err := db.Model(&models.User{}).Where("id = ?", personId).First(&models.User{}).Error
			if err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
		} else if tokenType == "owner" {
			var reqData GetGroupID
			if err := c.ShouldBindUri(&reqData); err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
			groupId := uuid.MustParse(reqData.ID)
			var role models.Role
			err := db.Model(&models.Role{}).Where("user_id = ? AND group_id = ? AND role = ?", personId, groupId, tokenType).First(&role).Error
			if err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
		} else if tokenType == "mentor" {
			var reqData GetMentorID
			if err := c.ShouldBindUri(&reqData); err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
			groupId := uuid.MustParse(reqData.GroupID)
			mentorId := uuid.MustParse(reqData.ID)
			var mentor models.Mentor
			err := db.Model(&models.Mentor{}).Where("user_id = ? AND group_id = ? AND mentor_id = ?", personId, groupId, mentorId).First(&mentor).Error
			if err != nil {
				httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
				c.Abort()
				return
			}
		}

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

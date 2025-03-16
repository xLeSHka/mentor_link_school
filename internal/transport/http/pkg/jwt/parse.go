package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func Parse(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get("personId")
	if !ok {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	uid, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	return uid, nil
}
func ParseGroupID(c *gin.Context) (uuid.UUID, error) {
	groupid := c.Param("groupID")
	if groupid == "" {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	return groupID, nil
}
func ParseUserID(c *gin.Context) (uuid.UUID, error) {
	userid := c.Param("userID")
	if userid == "" {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	userID, err := uuid.Parse(userid)
	if err != nil {
		return uuid.UUID{}, httpError.New(http.StatusUnauthorized, "Bad uuid")
	}
	return userID, nil
}

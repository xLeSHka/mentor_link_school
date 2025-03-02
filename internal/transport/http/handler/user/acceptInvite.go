package usersRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

func (h *Route) acceptedInvite(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	code := c.Query("code")
	if code == "" {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	ok, err := h.usersService.Invite(c.Request.Context(), code, personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	if !ok {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

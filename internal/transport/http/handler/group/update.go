package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

// @Summary Обновить код приглашения
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param id path string true "Group ID"
// @Success 200 {object} respUpdateCode
// @Failure 400 {object} httpError.HTTPError
// @Failure 401 {object} httpError.HTTPError
// @Router /api/groups/{groupID}/admin/inviteCode [post]
func (h *Route) updateInviteCode(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusBadRequest, "group not found").SendError(c)
		c.Abort()
		return
	}

	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusBadRequest, "group not found parse").SendError(c)
		c.Abort()
		return
	}

	code, err := h.groupService.UpdateInviteCode(c.Request.Context(), groupID, personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.JSON(http.StatusOK, respUpdateCode{
		Code: code,
	})
}

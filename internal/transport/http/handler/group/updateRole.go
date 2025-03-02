package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

// @Summary Обновить роль юзера
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/members/role [post]
// @Success 200
func (h *Route) updateRole(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	var req reqUpdateRole
	if err := h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, "Bad request").SendError(c)
		c.Abort()
		return
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	err = h.groupService.UpdateRole(c.Request.Context(), personID, groupID, userID, req.Role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	ws "github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Обновить роль юзера
// @Tags Roles
// @Accept json
// @Produce json
// @Router /api/groups/{id}/members/role [post]
// @Param id path string true "Group ID"
// @Param body body ReqUpdateRole true "body"
// @Param Authorization header string true "Bearer <token>"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 404 {object} httpError.HTTPError "Нет такого юзера"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) updateRole(c *gin.Context) {
	_, err := jwt.Parse(c)
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
	var req ReqUpdateRole
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
	role := &models.Role{
		GroupID: groupID,
		UserID:  userID,
		Role:    req.Role,
	}
	err = h.groupService.AddRole(c.Request.Context(), role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRole(userID, groupID, req.Role, h.producer, h.usersService, h.minioRepository, h.groupService)

	c.Writer.WriteHeader(http.StatusOK)
}

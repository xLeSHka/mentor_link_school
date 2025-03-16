package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"
)

// @Summary Удалить роль юзеру
// @Tags Roles
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/members/{userID}/role [delete]
// @Param groupID path string true "Group ID"
// @Param userID path string true "User ID"
// @Param body body ReqUpdateRole true "body"
// @Param Authorization header string true "Bearer <token>"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 404 {object} httpError.HTTPError "Нет такого юзера"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) removeRole(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	groupID, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	userID, err := jwt.ParseUserID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var req ReqUpdateRole
	if err := h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, "Bad request").SendError(c)
		c.Abort()
		return
	}

	role := &models.Role{
		GroupID: groupID,
		UserID:  userID,
		Role:    req.Role,
	}
	err = h.groupService.RemoveRole(c.Request.Context(), role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRole(userID, groupID, req.Role, h.producer, h.minioRepository, h.groupService)

	c.Writer.WriteHeader(http.StatusOK)
}

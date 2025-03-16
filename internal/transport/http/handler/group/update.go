package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Обновить код приглашения
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param groupID path string true "Group ID"
// @Success 200 {object} RespUpdateCode
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Router /api/groups/{groupID}/inviteCode [post]
// @Param Authorization header string true "Bearer <token>"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) updateInviteCode(c *gin.Context) {
	personID, err := jwt.Parse(c)
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

	code, err := h.groupService.UpdateInviteCode(c.Request.Context(), groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRole(personID, groupID, "owner", h.producer, h.minioRepository, h.groupService)
	c.JSON(http.StatusOK, RespUpdateCode{
		Code: code,
	})
}

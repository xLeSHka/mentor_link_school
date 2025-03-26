package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

// @Summary Создание организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqCreateGroupDto true "body"
// @Success 200 {object} RespCreateGroup
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Router /api/users/groups/create [post]
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) createGroup(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var reqData ReqCreateGroupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	group := &models.Group{
		ID:   uuid.New(),
		Name: reqData.Name,
	}

	inviteCode, err := h.groupService.Create(c.Request.Context(), group, personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRole(personId, group.ID, "owner", "", h.producer, h.minioRepository, h.groupService, h.usersService)
	c.JSON(http.StatusOK, RespCreateGroup{
		GroupID:    group.ID,
		InviteCode: inviteCode,
	})
}

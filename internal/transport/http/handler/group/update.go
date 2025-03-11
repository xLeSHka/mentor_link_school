package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Обновить код приглашения
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param id path string true "Group ID"
// @Success 200 {object} respUpdateCode
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Router /api/groups/{id}/inviteCode [post]
// @Param Authorization header string true "Bearer <token>"
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
	user, err := h.usersService.GetByID(c.Request.Context(), personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if user.AvatarURL != nil {
		avatrUrl, err := h.minioRepository.GetImage(*user.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		user.AvatarURL = &avatrUrl
	}
	group, err := h.usersService.GetGroupByID(c.Request.Context(), groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if group.AvatarURL != nil {
		avatrUrl, err := h.minioRepository.GetImage(*group.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		group.AvatarURL = &avatrUrl
	}
	role := "owner"
	mes := &ws.Message{
		Type:   "role",
		UserID: personID,
		Role: &ws.Role{
			Role:       role,
			GroupID:    groupID,
			GroupUrl:   group.AvatarURL,
			Name:       group.Name,
			InviteCode: group.InviteCode,
		},
	}
	go h.wsconn.WriteMessage(mes)
	c.JSON(http.StatusOK, respUpdateCode{
		Code: code,
	})
}

package groupsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
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
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Router /api/groups/{groupID}/inviteCode [post]
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
		Type:       "role",
		Role:       &role,
		GroupID:    &groupID,
		UserID:     &personID,
		GroupUrl:   group.AvatarURL,
		Name:       &group.Name,
		InviteCode: group.InviteCode,
	}
	go ws.WriteMessage(mes)
	c.JSON(http.StatusOK, respUpdateCode{
		Code: code,
	})
}

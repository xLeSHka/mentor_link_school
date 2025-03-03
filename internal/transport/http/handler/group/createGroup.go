package groupsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

// @Summary
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param body body reqCreateGroupDto true "body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} respCreateGroup
// @Router /api/groups/create [post]
func (h *Route) createGroup(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id")
		c.Abort()
		return
	}
	var reqData reqCreateGroupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	group := &models.Group{
		ID:   uuid.New(),
		Name: reqData.Name,
	}

	err = h.groupService.Create(c.Request.Context(), group, personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personId)
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
	group, err = h.usersService.GetGroupByID(c.Request.Context(), group.ID)
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
		Type:     "role",
		Role:     &role,
		GroupID:  &group.ID,
		UserID:   &personId,
		GroupUrl: group.AvatarURL,
		Name:     &group.Name,
	}
	if role == "owner" {
		mes.InviteCode = group.InviteCode
	}
	go ws.WriteMessage(mes)
	c.JSON(http.StatusOK, respCreateGroup{
		GroupID: group.ID,
	})
}

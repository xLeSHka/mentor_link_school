package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Редактирование организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{id}/edit [put]
// @Param id path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Param body body reqEditGroup true "body"
// @Failure 400 {object} httpError.HTTPError
// @Failure 401 {object} httpError.HTTPError
// @Success 200
func (h *Route) edit(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	var reqData reqEditGroup
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	toUpdate := make(map[string]any)
	toUpdate["name"] = reqData.Name
	group, err := h.groupService.Edit(c.Request.Context(), personID, groupID, toUpdate)
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
	c.Writer.WriteHeader(http.StatusOK)
}

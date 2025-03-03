package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

// @Summary Список участников организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Router /api/groups/{id}/members [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []resGetMember
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
func (h *Route) getMembers(c *gin.Context) {
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
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	members, err := h.groupService.GetMembers(c.Request.Context(), personID, groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*respGetMember, 0, len(members))
	for _, m := range members {
		if m.User.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.User.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			m.User.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMember(m))
	}
	c.JSON(http.StatusOK, resp)
}

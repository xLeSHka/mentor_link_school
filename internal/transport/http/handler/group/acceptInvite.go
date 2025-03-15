package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Присоединиться к организации по коду
// @Tags Roles
// @Accept json
// @Produce json
// @Router /api/groups/join/{code} [post]
// @Param code path string true "Invite code"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} respJoinGroup
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) acceptedInvite(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	code := c.Param("code")
	if code == "" {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	ok, err := h.usersService.Invite(c.Request.Context(), code, personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	if !ok {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	if h.producer != nil {
		user, err := h.usersService.GetByID(c.Request.Context(), personID)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		if user.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*user.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			user.AvatarURL = &avatarURL
		}
		group, err := h.usersService.GetGroupByInviteCode(c.Request.Context(), code)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		if group.AvatarURL != nil {
			groupAvatarURL, err := h.minioRepository.GetImage(*group.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			group.AvatarURL = &groupAvatarURL
		}
		go h.producer.Send(&ws.Message{
			Type:   "role",
			UserID: personID,
			Role: &ws.Role{
				Name:     user.Name,
				GroupID:  group.ID,
				GroupUrl: group.AvatarURL,
				Role:     "student",
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

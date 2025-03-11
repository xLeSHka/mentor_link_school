package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Присоединиться к организации по коду
// @Tags Groups
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
	groups, err := h.usersService.GetGroups(c.Request.Context(), personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*ws.RespGetGroupDto, 0, len(groups))
	for _, group := range groups {
		if group.Group.AvatarURL != nil {
			groupAvatarURL, err := h.minioRepository.GetImage(*group.Group.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			group.Group.AvatarURL = &groupAvatarURL
		}
		resp = append(resp, ws.MapGroup(group.Group, group.Role))
	}
	go h.wsconn.WriteMessage(&ws.Message{
		Type:   "user",
		UserID: personID,
		User: &ws.User{
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
			Telegram:  user.Telegram,
			BIO:       user.BIO,
			Groups:    resp,
		},
	})
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

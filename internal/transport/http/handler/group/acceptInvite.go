package groupsRoute

import (
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

// @Summary Присоединиться к группе по коду
// @Tags Groups
// @Accept json
// @Produce json
// @Router /groups/join/{code} [post]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} respJoinGrou
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
	group, err := h.usersService.GetGroupByInviteCode(c.Request.Context(), code)
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
	role := "student"
	mes := &ws.Message{
		Type:   "role",
		UserID: personID,
		Role: &ws.Role{
			Role:     role,
			GroupID:  group.ID,
			GroupUrl: group.AvatarURL,
			Name:     group.Name,
		},
	}
	go ws.WriteMessage(mes)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

// @Summary Редактирование организации
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /group/{id}/edit [put]
// @Security ApiKeyAuth
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
	err = h.usersService.Edit(c.Request.Context(), personID, groupID, toUpdate)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
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
	go ws.WriteMessage(&ws.Message{
		Type:   "user",
		UserID: personID,
		User: &ws.User{
			UserUrl:  user.AvatarURL,
			Telegram: user.Telegram,
			BIO:      user.BIO,
		},
	})
	c.Writer.WriteHeader(http.StatusOK)
}

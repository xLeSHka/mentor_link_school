package usersRoute

import (
	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

// @Summary Редактирование профиля
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /user/profile/edit [get]
// @Security ApiKeyAuth
// @Param body body reqEditUser true "body"
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
	var reqData reqEditUser
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	toUpdate := make(map[string]any)
	toUpdate["bio"] = reqData.BIO
	toUpdate["telegram"] = reqData.Telegram
	toUpdate["name"] = reqData.Name
	err = h.usersService.Edit(c.Request.Context(), personID, toUpdate)
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
		Type:     "user",
		UserID:   &personID,
		UserUrl:  user.AvatarURL,
		Telegram: &user.Telegram,
		BIO:      user.BIO,
	})
	c.Writer.WriteHeader(http.StatusOK)
}

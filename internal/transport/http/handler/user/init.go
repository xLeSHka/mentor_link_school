package usersRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// @Summary Получить инфу о себе
// @Schemes
// @Description Авторизация юзера
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/init [get]
// @Security ApiKeyAuth
// @Success 200 {object} resGetInitData
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) init(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
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
	groups, err := h.usersService.GetGroups(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*respGetGroupDto, 0, len(groups))
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
		resp = append(resp, mapGroup(group.Group, group.Role))
	}
	c.JSON(http.StatusOK, resGetInitData{
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		BIO:       user.BIO,
		Telegram:  &user.Telegram,
		Groups:    resp,
	})
}

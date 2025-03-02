package usersRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Получить инфу о себе
// @Schemes
// @Description Авторизация юзера
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/profile [get]
// @Success 200 {object} resGetProfile
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
func (h *Route) profile(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	user, err := h.usersService.GetByID(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}

	if user.AvatarURL != nil {
		avatarURL, err := h.minioRepository.GetImage(*user.AvatarURL)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		user.AvatarURL = &avatarURL
	}

	groups, err := h.usersService.GetGroups(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetGroupDto, 0, len(groups))
	for _, g := range groups {
		if g.Group.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*g.Group.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			g.Group.AvatarURL = &avatarURL
		}
		resp = append(resp, mapGroup(g, g.Role))
	}

	c.JSON(http.StatusOK, resGetProfile{
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		BIO:       user.BIO,
		Groups:    resp,
	})
}

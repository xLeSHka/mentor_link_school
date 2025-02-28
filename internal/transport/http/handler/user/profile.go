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
// @Router /api/user/me [get]
// @Success 200 {object} resGetProfile
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
func (h *Route) profile(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	user, err := h.usersService.GetByID(c.Request.Context(), personId)
	if err != nil {
		httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
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
	c.JSON(http.StatusOK, resGetProfile{
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		AvatarUrl: user.AvatarURL,
	})
}

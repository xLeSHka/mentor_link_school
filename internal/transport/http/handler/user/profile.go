package usersRoute

import (
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

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
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		user.AvatarURL = &avatarURL
	}
	c.JSON(http.StatusOK, resGetProfile{
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		BIO:       user.BIO,
		Telegram:  user.Telegram,
	})
}

func (h *Route) profileOther(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

	id := c.Param("id")
	if id == "" {
		httpError.New(http.StatusBadRequest, "Bad id").SendError(c)
		c.Abort()
		return
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		httpError.New(http.StatusBadRequest, "Bad id").SendError(c)
		c.Abort()
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), uuid)
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
	c.JSON(http.StatusOK, resGetProfile{
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		BIO:       user.BIO,
		Telegram:  user.Telegram,
	})
}

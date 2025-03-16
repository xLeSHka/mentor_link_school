package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// @Summary Получение профиля
// @Schemes
// @Description
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/profile/ [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} ResGetProfile
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) profile(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	err = avatar.GetUserAvatar(user, h.minioRepository)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, MapProfile(user))
}

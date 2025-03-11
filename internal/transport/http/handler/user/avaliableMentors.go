package usersRoute

import (
	"fmt"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

// @Summary Получение доступных менторов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/availableMentors [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []respGetMentor
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
func (h *Route) availableMentors(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

	mentors, err := h.usersService.GetMentors(c.Request.Context(), personId)
	if err != nil {
		fmt.Println(err)
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMentor, 0, len(mentors))
	for _, m := range mentors {
		if m.User.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.User.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			m.User.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMentor(m))
	}

	c.JSON(http.StatusOK, resp)
}

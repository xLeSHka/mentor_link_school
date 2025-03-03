package usersRoute

import (
	"fmt"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

// @Summary Получение доступных всех
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/availableMentors [get]
// @Security ApiKeyAuth
// @Success 200 {object} []respGetMentor
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
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

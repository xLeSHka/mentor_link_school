package usersRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Получение моих менторов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/mentors [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []respGetMyMentor
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
func (h *Route) getMyMentors(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

	mentors, err := h.usersService.GetMyMentors(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMyMentor, 0, len(mentors))
	for _, m := range mentors {
		if m.Mentor.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.Mentor.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			m.Mentor.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMyMentor(m))
	}

	c.JSON(http.StatusOK, resp)
}

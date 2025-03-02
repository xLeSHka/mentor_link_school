package usersRoute

import (
	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

// @Summary Получение моих менторов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/mentors [get]
// @Success 200 {object} []respGetMyMentor
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
				err.(*httpError.HTTPError).SendError(c)
				c.Abort()
				return
			}
			m.Mentor.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMyMentor(m))
	}

	c.JSON(http.StatusOK, resp)
}

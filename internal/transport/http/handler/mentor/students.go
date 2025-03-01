package mentorsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"
)

// @Summary Получить моих учеников
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/students [get]
// @Success 200 {object} []resGetProfile
func (h *Route) students(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	mentors, err := h.mentorService.GetStudents(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMyStudent, 0, len(mentors))
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
		resp = append(resp, mapMyMentor(m.Mentor))
	}

	c.JSON(http.StatusOK, resp)
}

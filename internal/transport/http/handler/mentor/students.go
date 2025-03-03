package mentorsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

// @Summary Получить моих учеников
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/students [get]
// @Success 200 {object} []respGetMyStudent
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) students(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

	students, err := h.mentorService.GetStudents(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMyStudent, 0, len(students))
	for _, m := range students {
		if m.Student.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.Student.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			m.Student.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMyStudent(m))
	}

	c.JSON(http.StatusOK, resp)
}

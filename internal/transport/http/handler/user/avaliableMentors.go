package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"
)

// @Summary Получение доступных всех
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/avaliableMentors [get]
// @Success 200 {object} []respGetMentor
func (h *Route) availableMentors(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	mentors, err := h.usersService.GetMentors(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMentor, 0, len(mentors))
	for _, m := range mentors {
		if m.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.AvatarURL)
			if err != nil {
				err.(*httpError.HTTPError).SendError(c)
				c.Abort()
				return
			}
			m.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMentor(m))
	}

	c.JSON(http.StatusOK, resp)
}

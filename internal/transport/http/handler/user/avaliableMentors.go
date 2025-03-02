package usersRoute

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
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
		fmt.Println(err)
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetMentor, 0, len(mentors))
	for _, m := range mentors {
		if m.User.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*m.User.AvatarURL)
			if err != nil {
				err.(*httpError.HTTPError).SendError(c)
				c.Abort()
				return
			}
			m.User.AvatarURL = &avatarURL
		}
		resp = append(resp, mapMentor(m))
	}

	c.JSON(http.StatusOK, resp)
}

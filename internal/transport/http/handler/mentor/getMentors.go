package mentorsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

// @Summary Получение списка менторов, которые доступны юзеру
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/group/{groupId}/mentors [get]
// @Success 200 {object} []models.User
func (h *Route) getMentors(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	_, err := h.mentorService.GetMentors(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

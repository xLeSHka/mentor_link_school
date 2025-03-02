package mentorsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

// @Summary Получить входящие запросы
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/requests [get]
// @Success 200 {object} []respGetRequest
func (h *Route) getRequests(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	mentors, err := h.mentorService.GetMyHelps(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetRequest, 0, len(mentors))
	for _, m := range mentors {
		resp = append(resp, mapRequest(m))
	}

	c.JSON(http.StatusOK, resp)
}

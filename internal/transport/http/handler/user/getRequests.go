package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"
)

// @Summary Получить список моих запросов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/requests [get]
// @Success 200 {object} []respGetHelp
func (h *Route) getRequests(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	mentors, err := h.usersService.GetMyHelps(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetHelp, 0, len(mentors))
	for _, m := range mentors {
		resp = append(resp, mapHelp(m))
	}

	c.JSON(http.StatusOK, resp)
}

package mentorsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

// @Summary Получить входящие запросы
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/requests [get]
// @Security ApiKeyAuth
// @Success 200 {object} []respGetRequest
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) getRequests(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

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

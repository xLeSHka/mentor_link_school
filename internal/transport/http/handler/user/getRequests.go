package usersRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Получить список моих запросов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/requests [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetHelp
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getRequests(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}

	mentors, err := h.usersService.GetMyHelps(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*RespGetHelp, 0, len(mentors))
	for _, m := range mentors {
		resp = append(resp, MapHelp(m))
	}

	c.JSON(http.StatusOK, resp)
}

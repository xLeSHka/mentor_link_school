package studentsRoute

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
// @Router /api/groups/{groupID}/students/requests [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetHelp
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getRequests(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	groupId, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	mentors, err := h.studentsService.GetMyHelps(c.Request.Context(), personId, groupId)
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

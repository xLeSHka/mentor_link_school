package mentorsRoute

import (
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

// @Summary Получить входящие запросы
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/mentors/requests [get]
// @Param groupID path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetRequest
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 404 {object} httpError.HTTPError "Нет такого запроса"
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
	mentors, err := h.mentorService.GetMyHelps(c.Request.Context(), personId, groupId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*RespGetRequest, 0, len(mentors))
	for _, m := range mentors {
		resp = append(resp, MapRequest(m))
	}

	c.JSON(http.StatusOK, resp)
}

package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Получить статистику организации
// @Tags Roles
// @Accept  json
// @Produce  json
// @Param groupID path string true "Group ID"
// @Router /api/groups/{groupID}/stat [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} RespStat
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getStat(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}

	groupID, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	stat, err := h.groupService.GetStat(c.Request.Context(), groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &RespStat{
		StudentsCount:        stat.StudentsCount,
		MentorsCount:         stat.MentorsCount,
		HelpRequestCount:     stat.HelpRequestCount,
		AcceptedRequestCount: stat.AcceptedRequestCount,
		RejectedRequestCount: stat.RejectedRequestCount,
		Conversion:           stat.Conversion,
	})

}

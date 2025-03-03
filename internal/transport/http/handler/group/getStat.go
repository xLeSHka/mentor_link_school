package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

// @Summary Получить статистику организации
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param id path string true "Group ID"
// @Router /groups/{GroupID}/stat [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} respStat
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
func (h *Route) getStat(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id")
		c.Abort()
		return
	}

	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	stat, err := h.groupService.GetStat(c.Request.Context(), personID, groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &respStat{
		StudentsCount:        stat.StudentsCount,
		MentorsCount:         stat.MentorsCount,
		TotalCount:           stat.TotalCount,
		HelpRequestCount:     stat.HelpRequestCount,
		AcceptedRequestCount: stat.AcceptedRequestCount,
		RejectedRequestCount: stat.RejectedRequestCount,
		Conversion:           stat.Conversion,
	})

}

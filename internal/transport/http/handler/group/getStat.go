package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

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

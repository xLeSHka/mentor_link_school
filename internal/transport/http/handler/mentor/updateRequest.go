package mentorsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

// @Summary Изменить состояние заявки
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/requests [post]
// @Param body body reqUpdateRequest true "body"
func (h *Route) updateRequest(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	var req reqUpdateRequest
	if err = h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, err.Error())
		return
	}
	var status string
	if req.Status {
		status = "accept"
	} else {
		status = "reject"
	}
	request := &models.HelpRequest{
		ID:       req.ID,
		MentorID: personId,
		Status:   status,
	}
	err = h.mentorService.UpdateRequest(c.Request.Context(), request)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

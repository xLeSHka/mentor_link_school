package mentorsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	personID := uuid.MustParse(c.MustGet("personId").(string))
	var req reqUpdateRequest
	if err := h.validator.ShouldBindJSON(c, &req); err != nil {
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
		MentorID: personID,
		Status:   status,
	}
	err := h.mentorService.UpdateRequest(c.Request.Context(), request)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

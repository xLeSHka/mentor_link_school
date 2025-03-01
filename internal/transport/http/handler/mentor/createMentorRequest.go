package mentorsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (h *Route) createMentorRequest(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))
	var goal GetMentorRequestDto
	if err := h.validator.ShouldBindJSON(c, &goal); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	var reqData GetGroupID
	if err := h.validator.ShouldBindUri(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	groupID := uuid.MustParse(reqData.ID)
	req := &models.CreateMentorRequest{
		UserID:  personId,
		GroupID: groupID,
		Goal:    goal.Goal,
	}

	err := h.mentorService.CreateMentorRequest(c.Request.Context(), req)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

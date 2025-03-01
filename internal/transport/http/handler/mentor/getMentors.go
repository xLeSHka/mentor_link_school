package mentorsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

func (h *Route) getMentors(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))

	mentors, err := h.mentorService.GetMentors(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

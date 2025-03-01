package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (h *Route) createGroup(c *gin.Context) {
	//personId := uuid.MustParse(c.MustGet("personId").(string))
	var reqData reqCreateGroupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	group := &models.Group{
		ID:   uuid.New(),
		Name: reqData.Name,
	}

	err := h.groupService.CreateGroup(c.Request.Context(), group, uuid.New()) //TODO refactor
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"group_id": group.ID,
	})
}

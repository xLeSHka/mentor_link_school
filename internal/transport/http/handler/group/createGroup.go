package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

// @Summary
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param body body reqCreateGroupDto true "body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} respCreateGroup
// @Router /api/group [post]
func (h *Route) createGroup(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))
	var reqData reqCreateGroupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	group := &models.Group{
		ID:   uuid.New(),
		Name: reqData.Name,
	}

	err := h.groupService.Create(c.Request.Context(), group, personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.JSON(http.StatusOK, respCreateGroup{
		GroupID: group.ID,
	})
}

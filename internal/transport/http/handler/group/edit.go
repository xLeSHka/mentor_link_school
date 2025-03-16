package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Редактирование организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/edit [patch]
// @Param groupID path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqEditGroup true "body"
// @Failure 400 {object} httpError.HTTPError
// @Failure 401 {object} httpError.HTTPError
// @Success 200
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) edit(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var reqData ReqEditGroup
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	groupID, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}

	_, err = h.groupService.Edit(c.Request.Context(), &models.Group{
		ID:   groupID,
		Name: reqData.Name,
	})
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRole(personID, groupID, "owner", h.producer, h.minioRepository, h.groupService)
	c.Writer.WriteHeader(http.StatusOK)
}

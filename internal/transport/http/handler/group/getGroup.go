package groupsRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (h *Route) getGroup(c *gin.Context) {
	//personId := uuid.MustParse(c.MustGet("personId").(string))
	var reqData GetGroupID
	if err := h.validator.ShouldBindUri(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	groupID := uuid.MustParse(reqData.ID)
	group := &models.Group{
		//UserID: personId,
		ID: groupID,
	}

	group, err := h.groupService.GetGroup(c.Request.Context(), group)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if group.AvatarURL != nil {
		avatarURL, err := h.minioRepository.GetImage(*group.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		group.AvatarURL = &avatarURL
	}
	c.JSON(http.StatusOK, mapGroup(group))
}

package groupsRoute

import (
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Список участников организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param groupID path string true "Group ID"
// @Router /api/groups/{groupID}/members [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetMember
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getMembers(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	groupID, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var req OffsetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	size := 10
	page := 0
	if req.Page != 0 {
		page = req.Page
	}
	if req.Size != 0 {
		size = req.Size
	}
	members, total, err := h.groupService.GetMembers(c.Request.Context(), groupID, page, size)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*RespGetMember, 0, len(members))
	for _, m := range members {
		err = avatar.GetUserAvatar(m, h.minioRepository)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		resp = append(resp, mapMember(m))
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(http.StatusOK, resp)
}

package usersRoute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"
)

// @Summary Получение организаций пользователя
// @Schemes
// @Description
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/groups/ [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []ResGetGroup
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getGroups(c *gin.Context) {
	personID, err := jwt.Parse(c)
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
	groups, total, err := h.usersService.GetGroups(c.Request.Context(), personID, page, size)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*ResGetGroup, 0, len(groups))
	for _, group := range groups {
		err = avatar.GetGroupAvatar(group.Group, h.minioRepository)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		resp = append(resp, MapGroup(group))
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(http.StatusOK, resp)
}

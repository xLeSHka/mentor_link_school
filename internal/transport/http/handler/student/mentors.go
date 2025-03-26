package studentsRoute

import (
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Получение моих менторов
// @Schemes
// @Tags Students
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/students/mentors [get]
// @Param groupID path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetMyMentor
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) getMyMentors(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	groupId, err := jwt.ParseGroupID(c)
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
	mentors, total, err := h.studentsService.GetMyMentors(c.Request.Context(), personId, groupId, page, size)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*RespGetMyMentor, 0, len(mentors))
	for _, m := range mentors {
		err = avatar.GetUserAvatar(m.Mentor, h.minioRepository)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			return
		}
		resp = append(resp, MapMyMentor(m))
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(http.StatusOK, resp)
}

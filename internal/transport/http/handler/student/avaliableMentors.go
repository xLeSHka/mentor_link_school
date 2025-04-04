package studentsRoute

import (
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

// @Summary Получение доступных менторов
// @Schemes
// @Tags Students
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/students/availableMentors [get]
// @Param groupID path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetMentor
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) availableMentors(c *gin.Context) {
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
	mentors, total, err := h.studentsService.GetMentors(c.Request.Context(), personId, groupId, page, size)
	if err != nil {
		fmt.Println(err)
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*RespGetMentor, 0, len(mentors))
	for _, m := range mentors {
		err = avatar.GetUserAvatar(m.User, h.minioRepository)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			return
		}
		resp = append(resp, MapMentor(m))
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(http.StatusOK, resp)
}

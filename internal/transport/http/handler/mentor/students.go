package mentorsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"log"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

// @Summary Получить моих учеников
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/mentors/students [get]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} []RespGetMyStudent
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) students(c *gin.Context) {
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
	students, err := h.mentorService.GetStudents(c.Request.Context(), personId, groupId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*RespGetMyStudent, 0, len(students))
	for _, m := range students {
		err = avatar.GetUserAvatar(m.Student, h.minioRepository)
		if err != nil {
			log.Println(err)
			return
		}
		resp = append(resp, MapMyStudent(m))
	}

	c.JSON(http.StatusOK, resp)
}

package mentorsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

// @Summary Изменить состояние заявки
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/mentors/requests [post]
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqUpdateRequest true "body"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 404 {object} httpError.HTTPError "Нет такого запроса"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) updateRequest(c *gin.Context) {
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
	var req ReqUpdateRequest
	if err = h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	var status string
	if *req.Status {
		status = "accepted"
	} else {
		status = "rejected"
	}
	request := &models.HelpRequest{
		ID:       req.ID,
		GroupID:  groupId,
		MentorID: personId,
		Status:   status,
	}
	err = h.mentorService.UpdateRequest(c.Request.Context(), request)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	request, err = h.studentsService.GetRequestByID(c.Request.Context(), req.ID, groupId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRequest(request.UserID, request.MentorID, request.ID, groupId, h.producer, h.usersService, h.minioRepository, h.studentsService)
	c.Writer.WriteHeader(http.StatusOK)
}

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
// @Router /api/mentors/requests [post]
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
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
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
		MentorID: personId,
		Status:   status,
	}
	err = h.mentorService.UpdateRequest(c.Request.Context(), request)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	request, err = h.userService.GetRequestByID(c.Request.Context(), req.ID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	go ws.SendRequest(request.UserID, request.MentorID, request.ID, h.producer, h.userService, h.minioRepository)
	c.Writer.WriteHeader(http.StatusOK)
}

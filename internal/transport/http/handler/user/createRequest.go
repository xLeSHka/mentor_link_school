package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Кинуть запрос ментору
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/requests [post]
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqCreateHelp true "body"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// Failure 409 {object} httpError.HTTPError "Запрос уже отправлен"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) createRequest(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	var reqData ReqCreateHelp
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	for _, r := range reqData.Requests {
		request := &models.HelpRequest{
			ID:       uuid.New(),
			UserID:   personId,
			MentorID: r.MentorID,
			GroupID:  r.GroupId,
			Goal:     reqData.Goal,
			Status:   "pending",
			BIO:      user.BIO,
		}
		err := h.usersService.CreateRequest(c.Request.Context(), request)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			return
		}
		go ws.SendRequest(personId, r.MentorID, request.ID, h.producer, h.usersService, h.minioRepository)
	}
	c.Writer.WriteHeader(http.StatusOK)
}

package mentorsRoute

import (
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

// @Summary Изменить состояние заявки
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/requests [post]
// @Param Authorization header string true "Bearer <token>"
// @Param body body reqUpdateRequest true "body"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 404 {object} httpError.HTTPError "Нет такого запроса"
func (h *Route) updateRequest(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	var req reqUpdateRequest
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
	updatedReq, err := h.userService.GetRequestByID(c.Request.Context(), req.ID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	student, err := h.userService.GetByID(c.Request.Context(), updatedReq.UserID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if student.AvatarURL != nil {
		avatarUrl, err := h.minioRepository.GetImage(*student.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		student.AvatarURL = &avatarUrl
	}
	mentor, err := h.userService.GetByID(c.Request.Context(), updatedReq.MentorID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if mentor.AvatarURL != nil {
		avatrUrl, err := h.minioRepository.GetImage(*mentor.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		mentor.AvatarURL = &avatrUrl
	}
	groupsIDs, err := h.userService.GetCommonGroups(updatedReq.UserID, updatedReq.MentorID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return

	}
	go ws.WriteMessage(&ws.Message{
		Type:   "request",
		UserID: student.ID,
		Request: &ws.Request{
			ID:              request.ID,
			StudentID:       student.ID,
			MentorID:        mentor.ID,
			MentorName:      mentor.Name,
			StudentName:     student.Name,
			MentorUrl:       mentor.AvatarURL,
			StudentUrl:      student.AvatarURL,
			StudentTelegram: student.Telegram,
			StudentBio:      student.BIO,
			MentorTelegram:  mentor.Telegram,
			MentorBio:       mentor.BIO,
			GroupIDs:        groupsIDs,
			Goal:            request.Goal,
			Status:          request.Status,
		},
	})
	c.Writer.WriteHeader(http.StatusOK)
}

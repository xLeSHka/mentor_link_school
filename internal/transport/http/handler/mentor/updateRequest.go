package mentorsRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"

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
// @Security ApiKeyAuth
// @Param body body reqUpdateRequest true "body"
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) updateRequest(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	var req reqUpdateRequest
	if err = h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, err.Error())
		return
	}
	var status string
	if req.Status {
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
		Type:        "request",
		ID:          &request.ID,
		UserID:      &student.ID,
		StudentID:   &student.ID,
		MentorID:    &mentor.ID,
		MentorName:  &mentor.Name,
		StudentName: &student.Name,
		MentorUrl:   mentor.AvatarURL,
		StudentUrl:  student.AvatarURL,
		GroupIDs:    &groupsIDs,
		Goal:        &request.Goal,
		BIO:         student.BIO,
		Status:      &request.Status,
	})
	c.Writer.WriteHeader(http.StatusOK)
}

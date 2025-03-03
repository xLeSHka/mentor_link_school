package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
)

// @Summary Кинуть запрос ментору
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/requests [post]
// @Security ApiKeyAuth
// @Param body body reqCreateHelp true "body"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) createRequest(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	var reqData reqCreateHelp
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error())
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
		student, err := h.usersService.GetByID(c.Request.Context(), personId)
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
		mentor, err := h.usersService.GetByID(c.Request.Context(), r.MentorID)
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
		groupsIDs, err := h.usersService.GetCommonGroups(personId, r.MentorID)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			return

		}
		go ws.WriteMessage(&ws.Message{
			Type:   "request",
			UserID: student.ID,
			Request: &ws.Request{
				ID:          request.ID,
				StudentID:   student.ID,
				MentorID:    mentor.ID,
				MentorName:  mentor.Name,
				StudentName: student.Name,
				MentorUrl:   mentor.AvatarURL,
				StudentUrl:  student.AvatarURL,
				GroupIDs:    groupsIDs,
				Goal:        request.Goal,
				Status:      request.Status,
			},
		})
	}
	c.Writer.WriteHeader(http.StatusOK)
}

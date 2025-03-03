package groupsRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

// @Summary Обновить роль юзера
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/members/role [post]
// @Param Authorization header string true "Bearer <token>"
// @Success 200
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 404 {object} httpError.HTTPError "Нет такого юзера"
func (h *Route) updateRole(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	var req reqUpdateRole
	if err := h.validator.ShouldBindJSON(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, "Bad request").SendError(c)
		c.Abort()
		return
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	err = h.groupService.UpdateRole(c.Request.Context(), personID, groupID, userID, req.Role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}

	user, err := h.usersService.GetByID(c.Request.Context(), userID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if user.AvatarURL != nil {
		avatrUrl, err := h.minioRepository.GetImage(*user.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		user.AvatarURL = &avatrUrl
	}
	group, err := h.usersService.GetGroupByID(c.Request.Context(), groupID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	if group.AvatarURL != nil {
		avatrUrl, err := h.minioRepository.GetImage(*group.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		group.AvatarURL = &avatrUrl
	}
	mes := &ws.Message{
		Type:   "role",
		UserID: userID,
		Role: &ws.Role{
			Role:     req.Role,
			GroupID:  groupID,
			GroupUrl: group.AvatarURL,
			Name:     group.Name,
		},
	}
	if req.Role == "owner" {
		mes.Role.InviteCode = group.InviteCode
	}
	go ws.WriteMessage(mes)
	c.Writer.WriteHeader(http.StatusOK)
}

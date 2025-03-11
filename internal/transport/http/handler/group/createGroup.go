package groupsRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

// @Summary Создание организации
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body reqCreateGroupDto true "body"
// @Success 200 {object} respCreateGroup
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Router /api/groups/create [post]
func (h *Route) createGroup(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id")
		c.Abort()
		return
	}
	var reqData reqCreateGroupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	group := &models.Group{
		ID:   uuid.New(),
		Name: reqData.Name,
	}

	err = h.groupService.Create(c.Request.Context(), group, personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	if user.AvatarURL != nil {
		avatarURL, err := h.minioRepository.GetImage(*user.AvatarURL)
		if err != nil {
			httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
			c.Abort()
			return
		}
		user.AvatarURL = &avatarURL
	}
	groups, err := h.usersService.GetGroups(c.Request.Context(), personId)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	resp := make([]*ws.RespGetGroupDto, 0, len(groups))
	for _, group := range groups {
		if group.Group.AvatarURL != nil {
			groupAvatarURL, err := h.minioRepository.GetImage(*group.Group.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			group.Group.AvatarURL = &groupAvatarURL
		}
		resp = append(resp, ws.MapGroup(group.Group, group.Role))
	}
	go h.wsconn.WriteMessage(&ws.Message{
		Type:   "user",
		UserID: personId,
		User: &ws.User{
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
			Telegram:  user.Telegram,
			BIO:       user.BIO,
			Groups:    resp,
		},
	})
	c.JSON(http.StatusOK, respCreateGroup{
		GroupID: group.ID,
	})
}

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
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqCreateGroupDto true "body"
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
	var reqData ReqCreateGroupDto
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
	if h.producer != nil {
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
		group, err = h.usersService.GetGroupByID(c.Request.Context(), group.ID)
		if err != nil {
			err.(*httpError.HTTPError).SendError(c)
			c.Abort()
			return
		}
		if group.AvatarURL != nil {
			groupAvatarURL, err := h.minioRepository.GetImage(*group.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			group.AvatarURL = &groupAvatarURL
		}
		go h.producer.Send(&ws.Message{
			Type:   "role",
			UserID: personId,
			Role: &ws.Role{
				Role:       "owner",
				Name:       user.Name,
				GroupID:    group.ID,
				GroupUrl:   group.AvatarURL,
				InviteCode: group.InviteCode,
			},
		})
	}
	c.JSON(http.StatusOK, respCreateGroup{
		GroupID: group.ID,
	})
}

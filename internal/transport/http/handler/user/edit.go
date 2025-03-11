package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Редактирование профиля
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/profile/edit [put]
// @Param Authorization header string true "Bearer <token>"
// @Param body body reqEditUser true "body"
// @Failure 400 {object} httpError.HTTPError
// @Failure 401 {object} httpError.HTTPError
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Success 200
func (h *Route) edit(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	var reqData reqEditUser
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	toUpdate := make(map[string]any)
	toUpdate["bio"] = reqData.BIO
	toUpdate["telegram"] = reqData.Telegram
	toUpdate["name"] = reqData.Name
	err = h.usersService.Edit(c.Request.Context(), personID, toUpdate)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personID)
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
	groups, err := h.usersService.GetGroups(c.Request.Context(), personID)
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
	go h.producer.Send(&ws.Message{
		Type:   "user",
		UserID: personID,
		User: &ws.User{
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
			Telegram:  user.Telegram,
			BIO:       user.BIO,
			Groups:    resp,
		},
	})
	c.Writer.WriteHeader(http.StatusOK)
}

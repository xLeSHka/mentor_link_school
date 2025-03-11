package usersRoute

import (
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"net/http"
	"path/filepath"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"

	"github.com/bachvtuan/mime2extension"
	"github.com/gin-gonic/gin"
)

// @Summary Загрузка аватарки для пользователя
// @Schemes
// @Description Загрузка аватарки. Возвращает ссылку на аватарку, которая действует 7 дней
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Изображение для загрузки"
// @Router /api/user/uploadAvatar [post]
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} respUploadAvatarDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
func (h *Route) uploadAvatar(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	temp, err := file.Open()
	if err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	ext := filepath.Ext(file.Filename)
	imagename := fmt.Sprintf("%s%s", personId.String(), ext)

	err, mimetype := mime2extension.Lookup(imagename)
	if err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	if mimetype != "image/jpeg" && mimetype != "image/png" && mimetype != "image/jpg" && mimetype != "image/webp" {
		httpError.New(http.StatusBadRequest, fmt.Errorf("Bad file type, need jpeg/png/jpg/webp, got %s", mimetype).Error()).SendError(c)
		c.Abort()
		return
	}
	f := &models.File{
		Filename: imagename,
		Size:     file.Size,
		File:     temp,
		Mimetype: mimetype,
	}
	imageURL, hErr := h.usersService.UploadImage(c.Request.Context(), f, personId)
	if hErr != nil {
		hErr.SendError(c)
		c.Abort()
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
	go h.producer.Send(&ws.Message{
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
	c.JSON(http.StatusOK, respUploadAvatarDto{Url: imageURL})
}

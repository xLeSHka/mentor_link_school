package groupsRoute

import (
	"fmt"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"
	"path/filepath"

	"github.com/bachvtuan/mime2extension"
	"github.com/gin-gonic/gin"
)

// @Summary Загрузка аватарки
// @Schemes
// @Description Загрузка аватарки. Возвращает ссылку на аватарку, которая действует 7 дней
// @Tags Groups
// @Accept multipart/form-data
// @Param image formData file true "Изображение для загрузки"
// @Produce json
// @Router /api/groups/{groupID}/uploadAvatar [post]
// @Param groupID path string true "Group ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} RespUploadAvatarDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 403 {object} httpError.HTTPError "Ошибка доступа"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Нет прав доступа"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) uploadAvatar(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	groupId, err := jwt.ParseGroupID(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
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
	imageURL, hErr := h.groupService.UploadImage(c.Request.Context(), f, groupId)
	if hErr != nil {
		hErr.SendError(c)
		c.Abort()
		return
	}
	go ws.SendRole(personId, groupId, "owner", h.producer, h.minioRepository, h.groupService)
	c.JSON(http.StatusOK, RespUploadAvatarDto{Url: imageURL})
}

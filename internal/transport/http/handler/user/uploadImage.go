package usersRoute

import (
	"fmt"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
	"path/filepath"

	"github.com/bachvtuan/mime2extension"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary загрузка аватарки
// @Schemes
// @Description Загрузка аватарки. Возвращает ссылку на аватарку, которая действует 7 дней
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/uploadAvatar [post]
// @Param body body reqUploadAvatarDto true "body"
// @Success 200 {object} respUploadAvatarDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
func (h *Route) uploadAvatar(c *gin.Context) {
	personId := uuid.MustParse(c.MustGet("personId").(string))
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

	c.JSON(http.StatusOK, respUploadAvatarDto{Url: imageURL})
}

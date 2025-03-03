package groupsRoute

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"net/http"
	"path/filepath"

	"github.com/bachvtuan/mime2extension"
	"github.com/gin-gonic/gin"
)

// @Summary загрузка аватарки
// @Schemes
// @Description Загрузка аватарки. Возвращает ссылку на аватарку, которая действует 7 дней
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/groups/{id}/uploadAvatar [post]
// @Success 200 {object} respUploadAvatarDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
func (h *Route) uploadAvatar(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Bad id").SendError(c)
		c.Abort()
		return
	}
	groupid := c.Param("id")
	if groupid == "" {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	groupID, err := uuid.Parse(groupid)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
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
	imageURL, hErr := h.groupService.UploadImage(c.Request.Context(), f, groupID, personId)
	if hErr != nil {
		hErr.SendError(c)
		c.Abort()
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personId)
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
	role := "owner"
	mes := &ws.Message{
		Type:   "role",
		UserID: personId,
		Role: &ws.Role{
			Role:       role,
			GroupID:    groupID,
			GroupUrl:   group.AvatarURL,
			Name:       group.Name,
			InviteCode: group.InviteCode,
		},
	}
	go ws.WriteMessage(mes)
	c.JSON(http.StatusOK, respUploadAvatarDto{Url: imageURL})
}

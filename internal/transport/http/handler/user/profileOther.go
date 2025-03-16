package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Получение чужого профиля
// @Schemes
// @Description
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/profile/{profileID} [get]
// @Param profileID path string true "Profile ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} RespOtherProfile
// @Failure 400 {object} httpError.HTTPError "Невалидный запрос"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) profileOther(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var reqData ReqOtherProfileDto
	if err := h.validator.ShouldBindUri(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, "Bad Request, need id to get other profile").SendError(c)
		c.Abort()
		return
	}
	profileID, err := uuid.Parse(*reqData.ProfileID)
	if err != nil {
		httpError.New(http.StatusUnauthorized, "Header not found").SendError(c)
		c.Abort()
		return
	}
	user, err := h.usersService.GetByID(c.Request.Context(), personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	profile, err := h.usersService.GetByID(c.Request.Context(), profileID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	err = avatar.GetUserAvatar(user, h.minioRepository)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, MapOtherProfile(profile))
}

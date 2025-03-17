package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"
)

// @Summary Редактирование профиля
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/profile/edit [patch]
// @Param Authorization header string true "Bearer <token>"
// @Param body body ReqEditUser true "body"
// @Failure 400 {object} httpError.HTTPError
// @Failure 401 {object} httpError.HTTPError
// @Failure 404 {object} httpError.HTTPError "Нет такого пользователя"
// @Failure 409 {object} httpError.HTTPError "Пользователь с таким telegram уже зарегистрирован"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
// @Success 200
func (h *Route) edit(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var reqData ReqEditUser
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}
	if reqData.BIO == nil && reqData.Name == nil && reqData.Telegram == nil {
		httpError.New(http.StatusBadRequest, "you are not update user!").SendError(c)
		c.Abort()
		return
	}
	user := &models.User{}
	if reqData.BIO != nil {
		user.BIO = reqData.BIO
	}
	if reqData.Name != nil {
		user.Name = *reqData.Name
	}
	if reqData.Telegram != nil {
		user.Telegram = *reqData.Telegram
	}

	_, err = h.usersService.Edit(c.Request.Context(), personID, user)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	go ws.SendUser(personID, h.producer, h.usersService, h.minioRepository)
	c.Writer.WriteHeader(http.StatusOK)
}

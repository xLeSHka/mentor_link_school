package usersRoute

import (
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"

	"github.com/gin-gonic/gin"
)

// @Summary Аунтефикация
// @Schemes
// @Description Аунтефикация юзера. Возвращает токен, который в дальнейшем нужно передавать в заголовке "Authorization" в формате "Bearer <токен>".
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/auth/login [post]
// @Param body body ReqLoginDto true "body"
// @Success 200 {object} RespLoginDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 403 {object} httpError.HTTPError "Пользователь заблокирован"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) login(c *gin.Context) {
	var reqData ReqLoginDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	token, err := h.usersService.Login(c.Request.Context(), *reqData.Telegram, *reqData.Password)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, RespLoginDto{
		Token: token,
	})
}

package usersRoute

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Аунтефикация
// @Schemes
// @Description Аунтефикация юзера. Возвращает токен, который в дальнейшем нужно передавать в заголовке "Authorization" в формате "Bearer <токен>". Кусок "Bearer " нужно добавлять самому. Это фиксированное слово, которое ставят перед токеном зачем-то.
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/auth/sign-in [post]
// @Param body body reqSigninDto true "body"
// @Success 200 {object} resSigninDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Неверный email или пароль"
func (h *Route) signin(c *gin.Context) {
	var reqData reqSigninDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	_, token, err := h.usersService.Login(c.Request.Context(), reqData.Email, reqData.Password)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.JSON(http.StatusOK, resSigninDto{
		Token: token,
	})
}

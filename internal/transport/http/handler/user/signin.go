package usersRoute

import (
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"

	"github.com/gin-gonic/gin"
)

// @Summary Аунтефикация
// @Schemes
// @Description Аунтефикация юзера. Возвращает токен, который в дальнейшем нужно передавать в заголовке "Authorization" в формате "Bearer <токен>". Кусок "Bearer " нужно добавлять самому. Это фиксированное слово, которое ставят перед токеном зачем-то.
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/auth/sign-in [post]
// @Param body body reqLoginDto true "body"
// @Success 200 {object} respLoginDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Неверный email или пароль"
func (h *Route) login(c *gin.Context) {
	var reqData reqLoginDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	_, token, err := h.usersService.Login(c.Request.Context(), reqData.Name)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	c.JSON(http.StatusOK, respLoginDto{
		Token: token,
	})
}

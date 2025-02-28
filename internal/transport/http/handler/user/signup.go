package usersRoute

import (
	"net/http"
	"prodapp/internal/app/httpError"
	"prodapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Регистрация
// @Schemes
// @Description Регистрация юзера. Возвращает токен, который в дальнейшем нужно передавать в заголовке "Authorization" в формате "Bearer <токен>". Кусок "Bearer " нужно добавлять самому. Это фиксированное слово, которое ставят перед токеном зачем-то.
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/user/auth/sign-up [post]
// @Param body body reqSignupDto true "body"
// @Success 200 {object} resSignupDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
func (h *Route) signup(c *gin.Context) {
	var reqData reqSignupDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		return
	}

	company := &models.User{
		ID:        uuid.New(),
		Name:      reqData.Name,
		Surname:   reqData.Surname,
		Email:     reqData.Email,
		AvatarURL: reqData.AvatarUrl,
		Password:  []byte(reqData.Password),
		Age:       reqData.Age,
	}

	token, err := h.usersService.Create(c.Request.Context(), company)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}

	c.JSON(http.StatusOK, resSignupDto{
		Token: token,
	})
}

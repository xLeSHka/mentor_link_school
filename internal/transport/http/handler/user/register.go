package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/password"
	"net/http"
)

// @Summary регистрация
// @Schemes
// @Description регистрация юзера. Возвращает токен, который в дальнейшем нужно передавать в заголовке "Authorization" в формате "Bearer <токен>".
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/auth/register [post]
// @Param body body ReqRegisterDto true "body"
// @Success 200 {object} RespRegisterDto
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 409 {object} httpError.HTTPError "Пользователь с таким email уже зарегистрирован"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) register(c *gin.Context) {
	var reqData ReqRegisterDto
	if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	encryprted, err := password.Encrypt([]byte(*reqData.Password), h.cryptoKey)
	if err != nil {
		httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
		c.Abort()
		return
	}
	user := &models.User{
		ID:       uuid.New(),
		Name:     *reqData.Name,
		Telegram: *reqData.Telegram,
		Password: encryprted,
	}
	token, err := h.usersService.Register(c.Request.Context(), user)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, RespLoginDto{
		Token: token,
	})
}

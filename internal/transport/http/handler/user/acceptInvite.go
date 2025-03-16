package usersRoute

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
)

// @Summary Присоединиться к организации по коду
// @Tags Users
// @Accept json
// @Produce json
// @Router /api/users/join/{code} [post]
// @Param code path string true "Invite code"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} RespJoinGroup
// @Failure 400 {object} httpError.HTTPError "Ошибка валидации"
// @Failure 401 {object} httpError.HTTPError "Ошибка авторизации"
// @Failure 500 {object} httpError.HTTPError "Что-то пошло не так"
func (h *Route) acceptedInvite(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	code := c.Param("code")
	if code == "" {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	ok, err := h.usersService.Invite(c.Request.Context(), code, personID)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	if !ok {
		httpError.New(http.StatusBadRequest, "Code not found").SendError(c)
		c.Abort()
		return
	}
	group, err := h.usersService.GetGroupByInviteCode(c.Request.Context(), code)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	go ws.SendRole(personID, group.ID, "student", h.producer, h.minioRepository, h.groupService)
	c.JSON(http.StatusOK, RespJoinGroup{Status: "ok"})
}

package usersRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

func (h *Route) getGroups(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	role := c.Query("role")
	if role == "" {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}

	groups, err := h.usersService.GetGroups(c.Request.Context(), personId, role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetGroupDto, 0, len(groups))
	for _, g := range groups {
		if g.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*g.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			g.AvatarURL = &avatarURL
		}
		resp = append(resp, mapGroup(g, req.Role))
	}
	c.JSON(http.StatusOK, resp)
}

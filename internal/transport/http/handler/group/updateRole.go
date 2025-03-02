package groupsRoute

import "github.com/gin-gonic/gin"

// @Summary Обновить роль юзера
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/admin/members/{memberID}/role [post]
// @Success 200 {object} respGetGroupDto
func (h *Route) updateRole(c *gin.Context) {}

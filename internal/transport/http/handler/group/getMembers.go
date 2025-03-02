package groupsRoute

import "github.com/gin-gonic/gin"

// @Summary Список участников группы
// @Schemes
// @Tags Groups
// @Accept json
// @Produce json
// @Router /api/groups/{groupID}/admin/members [post]
// @Success 200 {object} []resGetMember
func (h *Route) getMembers(c *gin.Context) {}

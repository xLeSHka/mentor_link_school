package mentorsRoute

import "github.com/gin-gonic/gin"

// @Summary Получить входящие запросы
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/requests [get]
// @Success 200 {object} []respGetRequest
func (h *Route) getRequests(c *gin.Context) {}

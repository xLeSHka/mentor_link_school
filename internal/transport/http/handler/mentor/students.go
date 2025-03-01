package mentorsRoute

import "github.com/gin-gonic/gin"

// @Summary Получить моих учеников
// @Schemes
// @Tags Mentors
// @Accept json
// @Produce json
// @Router /api/mentors/students [get]
// @Success 200 {object} []resGetProfile
func (h *Route) students(c *gin.Context) {}

package publicRoute

import (
	"github.com/gin-gonic/gin"
)

func (r *Route) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}

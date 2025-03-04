package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
)

func (r *Route) ping(c *gin.Context) {
	go ws.WriteMessage(&ws.Message{
		Type:    "request",
		UserID:  uuid.New(),
		Request: &ws.Request{},
	})

	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}

package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
func (r *Route) Websocket(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	for {
		time.Sleep(10 * time.Second)
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello world"))
		if err != nil {
			log.Println("write:", err)
			ws.Close()
			return
		}
	}
}

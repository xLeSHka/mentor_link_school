package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/pkg/jwt"
	"log"
	"net/http"
)

func (p *WebSocket) WriteMessage(message *Message) {
	p.Broadcast <- message
}

func (p *WebSocket) WsHandler(c *gin.Context) {
	//c.Writer.Header().Set("Connection", "Upgrade")
	//c.Writer.Header().Set("Upgrade", "websocket")
	println("wsHandler")
	personID, err := jwt.Parse(c)
	if err != nil {
		fmt.Println(err)
		httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
		c.Abort()
		return
	}
	// use default options
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected")
	// register client
	p.Clients[personID] = ws
}

func (p *WebSocket) Echo() {

	for {
		val := <-p.Broadcast
		jsonData, err := json.Marshal(val)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Message: " + string(jsonData))
		if val.Type == "request" {
			client, ok := p.Clients[val.Request.MentorID]
			if ok {
				err := client.WriteMessage(websocket.BinaryMessage, jsonData)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(p.Clients, val.Request.MentorID)
				}
			}
		}
		// send to every client that is currently connected
		client, ok := p.Clients[val.UserID]
		if ok {
			err := client.WriteMessage(websocket.BinaryMessage, jsonData)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(p.Clients, val.UserID)
			}
		}
	}
}
func (p *WebSocket) wsHandler(c *gin.Context) {
	//c.Writer.Header().Set("Connection", "Upgrade")
	//c.Writer.Header().Set("Upgrade", "websocket")
	println("wsHandler")
	personID, err := jwt.Parse(c)
	if err != nil {
		fmt.Println(err)
		httpError.New(http.StatusUnauthorized, err.Error()).SendError(c)
		c.Abort()
		return
	}
	// use default options
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected")
	// register client
	p.Clients[personID] = ws

}

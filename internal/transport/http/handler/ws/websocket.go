package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
)

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type Message struct {
	Type   string     `json:"type"`              // user | role | request | group
	UserID *uuid.UUID `json:"user_id,omitempty"` // есть везде
	BIO    *string    `json:"bio,omitempty"`     // в user и в request

	//общее для role и group. Является group сообщением по факту
	GroupID    *uuid.UUID `json:"group_id,omitempty"`
	Name       *string    `json:"name,omitempty"`
	GroupUrl   *string    `json:"group_url,omitempty"`
	InviteCode *string    `json:"invite_code,omitempty"`

	// сообщение о обновлении данных пользователя
	UserUrl  *string `json:"user_url,omitempty"`
	Telegram *string `json:"telegram,omitempty"`

	// сообщение о изменении ролей пользователя
	Role *string `json:"role,omitempty"`

	//сообщение о создании/изменении запроса на получение ментора
	ID          *uuid.UUID   `json:"id,omitempty"`
	StudentID   *uuid.UUID   `json:"student_id,omitempty"`
	MentorID    *uuid.UUID   `json:"mentor_id,omitempty"`
	MentorName  *string      `json:"mentor_name,omitempty"`
	StudentName *string      `json:"student_name,omitempty"`
	MentorUrl   *string      `json:"mentor_url,omitempty"`
	StudentUrl  *string      `json:"student_url,omitempty"`
	GroupIDs    *[]uuid.UUID `json:"group_ids,omitempty"`
	Goal        *string      `json:"goal,omitempty"`
	Status      *string      `json:"status,omitempty"`

	//

}

var clients = make(map[uuid.UUID]*websocket.Conn)
var broadcast = make(chan *Message)

func WriteMessage(message *Message) {
	broadcast <- message
}

func WsHandler(c *gin.Context) {
	personID, err := jwt.Parse(c)
	if err != nil {
		httpError.New(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected")
	// register client
	clients[personID] = ws
}

func Echo() {
	for {
		val := <-broadcast
		jsonData, err := json.Marshal(val)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Message: " + string(jsonData))
		if val.Type == "request" {
			client, ok := clients[*val.MentorID]
			if ok {
				err := client.WriteMessage(websocket.BinaryMessage, jsonData)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, *val.MentorID)
				}
			}
		}
		// send to every client that is currently connected
		client, ok := clients[*val.UserID]
		if ok {
			err := client.WriteMessage(websocket.BinaryMessage, jsonData)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, *val.UserID)
			}
		}
	}
}

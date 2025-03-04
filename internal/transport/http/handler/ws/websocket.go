package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
type Role struct {
	Role       string    `json:"role"`
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	GroupUrl   *string   `json:"group_url,omitempty"`
	InviteCode *string   `json:"invite_code,omitempty"`
}
type User struct {
	UserUrl  *string `json:"user_url,omitempty"`
	Telegram string  `json:"telegram"`
	BIO      *string `json:"bio,omitempty"`
}
type Request struct {
	ID              uuid.UUID   `json:"id"`
	StudentID       uuid.UUID   `json:"student_id"`
	MentorID        uuid.UUID   `json:"mentor_id"`
	MentorName      string      `json:"mentor_name"`
	StudentName     string      `json:"student_name"`
	MentorUrl       *string     `json:"mentor_url,omitempty"`
	StudentUrl      *string     `json:"student_url,omitempty"`
	StudentBio      *string     `json:"student_bio,omitempty"`
	MentorBio       *string     `json:"mentor_bio,omitempty"`
	StudentTelegram string      `json:"student_telegram"`
	MentorTelegram  string      `json:"mentor_telegram"`
	GroupIDs        []uuid.UUID `json:"group_ids"`
	Goal            string      `json:"goal"`
	Status          string      `json:"status"`
}
type Message struct {
	Type    string    `json:"type"`
	UserID  uuid.UUID `json:"user_id"`
	Role    *Role     `json:"role,omitempty"`
	User    *User     `json:"user,omitempty"`
	Request *Request  `json:"request,omitempty"`
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
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := ws.WriteMessage(websocket.TextMessage, []byte("hello"))
			if err != nil {
				log.Println(err)
				ws.Close()
				//delete(clients, personID)
				return
			}
		}

	}()
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
			client, ok := clients[val.Request.MentorID]
			if ok {
				err := client.WriteMessage(websocket.BinaryMessage, jsonData)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, val.Request.MentorID)
				}
			}
		}
		// send to every client that is currently connected
		client, ok := clients[val.UserID]
		if ok {
			err := client.WriteMessage(websocket.BinaryMessage, jsonData)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, val.UserID)
			}
		}
	}
}

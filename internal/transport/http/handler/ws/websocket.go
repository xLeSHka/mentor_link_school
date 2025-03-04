package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"log"
	"net/http"
)

type Role struct {
	Role       string    `json:"role"`
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	GroupUrl   *string   `json:"group_url,omitempty"`
	InviteCode *string   `json:"invite_code,omitempty"`
}
type RespGetGroupDto struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	AvatarUrl  *string `json:"avatar_url,omitempty"`
	InviteCode *string `json:"invite_code,omitempty"`
	Role       string  `json:"role"`
}

func MapGroup(group *models.Group, role string) *RespGetGroupDto {
	resp := &RespGetGroupDto{
		Name:      group.Name,
		ID:        group.ID.String(),
		AvatarUrl: group.AvatarURL,
		Role:      role,
	}
	if role == "owner" {
		resp.InviteCode = group.InviteCode
	}
	return resp
}

type User struct {
	Name      string             `json:"name"`
	AvatarUrl *string            `json:"avatar_url,omitempty"`
	Telegram  string             `json:"telegram"`
	BIO       *string            `json:"bio,omitempty"`
	Groups    []*RespGetGroupDto `json:"groups"`
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
type WebSocket struct {
	Conn      *websocket.Conn
	Clients   map[uuid.UUID]*websocket.Conn
	Broadcast chan *Message
}

func New() *WebSocket {
	var clients = make(map[uuid.UUID]*websocket.Conn)
	var broadcast = make(chan *Message)

	return &WebSocket{
		Conn:      nil,
		Clients:   clients,
		Broadcast: broadcast,
	}
}

func (p *WebSocket) WriteMessage(message *Message) {

	p.Broadcast <- message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
	//go func() {
	//	ticker := time.NewTicker(3 * time.Second)
	//	defer ticker.Stop()
	//	for range ticker.C {
	//		err := ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	//		if err != nil {
	//			log.Println(err)
	//			//ws.Close()
	//			//delete(clients, personID)
	//			return
	//		}
	//	}
	//
	//}()
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

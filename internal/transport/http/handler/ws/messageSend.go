package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

func (p *WebSocket) WriteMessage() {
	for {
		select {
		case data, ok := <-p.Consumer.Messages:
			if !ok {
				log.Println("Message channel closed")
				return
			}
			var m Message
			err := json.Unmarshal(data, &m)
			if err != nil {
				log.Println(err)
				continue
			}
			p.Broadcast <- &m
		}
	}
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

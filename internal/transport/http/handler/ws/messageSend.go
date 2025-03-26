package ws

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		switch val.Type {
		case "role":
			switch val.Role.Action {
			case "add":
				id, err := p.UsersService.GetTelegramID(context.Background(), val.UserID)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = p.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Вам добавили роль %s в организации %s", Roles(val.Role.Role), val.Role.Name)))
				if err != nil {
					log.Println(err)
				}
			case "remove":
				id, err := p.UsersService.GetTelegramID(context.Background(), val.UserID)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = p.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Вам удалили роль %s в организации %s", Roles(val.Role.Role), val.Role.Name)))
				if err != nil {
					log.Println(err)
				}
			}
		case "request":
			switch val.Request.Status {
			case "pending":
				id, err := p.UsersService.GetTelegramID(context.Background(), val.Request.MentorID)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = p.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("%s отправил вам запрос на менторство с целью %s", val.Request.StudentName, val.Request.Goal)))
				if err != nil {
					log.Println(err)
				}
			case "accepted":
				id, err := p.UsersService.GetTelegramID(context.Background(), val.Request.MentorID)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = p.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Ваш запрос ментору %s, с целью %s был принят 🤩", val.Request.MentorName, val.Request.Goal)))
				if err != nil {
					log.Println(err)
				}
			case "rejected":
				id, err := p.UsersService.GetTelegramID(context.Background(), val.Request.MentorID)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = p.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Ваш запрос ментору %s, с целью %s был отклонён 😢", val.Request.MentorName, val.Request.Goal)))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
func Roles(role string) string {
	switch role {
	case "student":
		return "студента"
	case "mentor":
		return "ментора"
	}
	return ""
}

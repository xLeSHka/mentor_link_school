package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"log"
	"net/http"
)

func SendRequest(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = SendRequest // Set self as current Action
	data := userDatas[stack.ChatID]
	data.Req.ID = uuid.New()
	data.Req.UserID = data.User.ID
	data.Req.MentorID = data.Profile.ID
	data.Req.GroupID = data.Group.ID
	data.Req.BIO = data.User.BIO
	data.Req.Status = "pending"
	if stack.IsPrint {
		log.Println(data.LastMes)
		stack.IsPrint = false
		if data.LastMes == -1 {
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПожалуйста введите цель менторства!", SendReqMenuTemplate))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Req = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			data.LastMes = sended.MessageID
			return stack
		} else {
			// Print UI
			msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\nПожалуйста введите цель менторства!", SendReqMenuTemplate))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			_, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Req = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		}
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					log.Println("Back", data.LastMes)
					userDatas[stack.ChatID].Req = nil
					return ReturnOnParent(stack)
				}
			}
		}
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			data.LastMes = -1
			data.Req.Goal = msgText
			err := stack.Bot.StudentService.CreateRequest(context.Background(), data.Req)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Req = nil
				data.LastMes = -1
				if err.(*httpError.HTTPError).StatusCode == http.StatusBadRequest {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой ментор не найден!🤨🔎", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nВы уже отправили запрос этому ментору с этой целью!🚫", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				}
				return ReturnOnParent(stack)
			}
			id, err := stack.Bot.UsersService.GetTelegramID(context.Background(), data.Req.MentorID)
			if err != nil {
				log.Println(err)
				return ReturnOnParent(stack)
			}
			_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("%s отправил вам запрос на менторство с целью %s", data.User.Name, data.Req.Goal)))
			if err != nil {
				log.Println(err)
			}
			return Request(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Update:  nil,
				Data:    "Created",
			})

		}
	}
	return stack
}

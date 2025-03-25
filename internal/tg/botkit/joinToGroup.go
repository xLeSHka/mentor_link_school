package botkit

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func JoinToGroup(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = JoinToGroup // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		if stack.LastMes != -1 {
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\nПожалуйста введите пригласительный код!", JoinMenuTemplate))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			_, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		} else {
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПожалуйста введите пригласительный код!", JoinMenuTemplate))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			stack.LastMes = sended.MessageID
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
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
			default:
				return stack
			}
		}
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			group, err := stack.Bot.UserRepository.GetGroupByInviteCode(context.Background(), msgText)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				if errors.Is(err, gorm.ErrRecordNotFound) {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Организация не найдена!🤨🔎 Введите другой инвайт код!"))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					sended, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.LastMes = sended.MessageID
					return stack
				}
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate))
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				stack.LastMes = sended.MessageID
				return ReturnOnParent(stack)
			}
			group.InviteCode = nil
			userDatas[stack.ChatID].Group = group
			_, err = stack.Bot.UsersService.Invite(context.Background(), msgText, data.User.ID)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil

				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nОрганизация с таким инвайт кодом не найдена!🤨🔎", ErrorMenuTemplate))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					sended, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.LastMes = sended.MessageID
					return stack
				} else if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nВы уже состоите в этой организации!🚫", ErrorMenuTemplate))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					sended, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.LastMes = sended.MessageID
				} else {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					sended, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.LastMes = sended.MessageID
				}
				return ReturnOnParent(stack)
			}

			return Group(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Update:  nil,
				LastMes: -1,
				Data:    "Created",
			})

		}
	}
	return stack
}

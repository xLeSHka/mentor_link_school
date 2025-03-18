package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
)

func CreateGroup(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = CreateGroup // Set self as current Action
	data := userDatas[stack.ChatID]
	data.Group.ID = uuid.New()
	if stack.IsPrint {
		stack.IsPrint = false
		// Print UI
		msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя организации!", CreateGroupMenuTemplate, CreateGroupTextTemplate(data.Group.ID, "____", "____")))
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

			if len(msgText) > 100 || len(msgText) < 1 {
				minL, maxL := "✔️", "✔️"
				if len(msgText) > 100 {
					maxL = "❌"
				}
				if len(msgText) < 1 {
					minL = "❌"
				}
				_, err := stack.Bot.Api.Send(tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\nНевалидне имя!\nМинимальная длина: %s\nМаксимальная длина: %s\nПожалуйста введите другое имя!", CreateGroupMenuTemplate, minL, maxL)))
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				return stack
			}
			data.Group.Name = msgText
			inviteCode, err := stack.Bot.GroupService.Create(context.Background(), data.Group, data.User.ID)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				_, err := stack.Bot.Api.Send(tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			data.Group.InviteCode = &inviteCode

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

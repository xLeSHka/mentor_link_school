package botkit

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Регистрация"),
		tgbotapi.NewKeyboardButton("Авторизация"),
	),
)

func MainMenu(stack CallStack) CallStack {
	stack.Action = MainMenu
	data := userDatas[stack.ChatID]
	if data == nil {
		if stack.Update != nil {
			if stack.Update.Message != nil {
				switch stack.Update.Message.Text {
				case "Register":
					userDatas[stack.ChatID] = &Data{}
					return RegisterName(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						IsPrint: true,
						Parent:  &stack,
						Update:  nil,
					})
				case "Login":
					userDatas[stack.ChatID] = &Data{}
					return Chop(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						IsPrint: true,
					})
				}
			}
		}
		text := fmt.Sprintf("[Главное меню]\n\nПожалуйста выберите одно из действий 🙏")
		msg := tgbotapi.NewMessage(stack.ChatID, text)
		msg.ReplyMarkup = mainMenuKeyboard
		_, err := stack.Bot.Api.Request(msg)
		if err != nil {
			log.Println(err)
		}
		return stack
	}
	return Chop(CallStack{
		ChatID:  stack.ChatID,
		Bot:     stack.Bot,
		IsPrint: true,
	})
}

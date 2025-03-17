package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
)

func MainMenuKeyboard(bot *Bot, telegram string) (tgbotapi.ReplyKeyboardMarkup, error) {
	_, err := bot.UsersService.GetByTelegram(context.Background(), telegram)
	if err != nil && err.(*httpError.HTTPError).StatusCode == 500 {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	if err != nil && err.(*httpError.HTTPError).StatusCode == 404 {
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Регистрация"),
			),
		), nil
	}
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Авторизация"),
		),
	), nil
}

func MainMenu(stack CallStack) CallStack {
	stack.Action = MainMenu
	data := userDatas[stack.ChatID]
	if data.User == nil {
		if stack.Update != nil {
			if stack.Update.Message != nil {
				switch stack.Update.Message.Text {
				case "Регистрация":
					userDatas[stack.ChatID] = &Data{User: &models.User{}}
					return RegisterName(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						IsPrint: true,
						Parent:  &stack,
						Update:  nil,
						Data:    stack.Data,
					})
				case "Авторизация":
					userDatas[stack.ChatID] = &Data{User: &models.User{}}
					return LoginPassword(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						Parent:  &stack,
						IsPrint: true,
						Update:  nil,
						Data:    stack.Data,
					})
				}
			}
		}

		keyboard, err := MainMenuKeyboard(stack.Bot, stack.Data)
		if err != nil {
			log.Println(err)
			return stack
		}
		text := fmt.Sprintf("%s\n\n%s", MainMenuTemplate, MainMenuTextTemplate)
		msg := tgbotapi.NewMessage(stack.ChatID, text)
		msg.ReplyMarkup = keyboard
		_, err = stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			return stack
		}
		return stack
	}
	mainMenuKeyboard := tgbotapi.NewRemoveKeyboard(true)
	msg := tgbotapi.NewMessage(stack.ChatID, "")
	msg.ReplyMarkup = mainMenuKeyboard
	_, _ = stack.Bot.Api.Send(msg)

	return AuthedMenu(CallStack{
		ChatID:  stack.ChatID,
		Bot:     stack.Bot,
		IsPrint: true,
		Parent:  &stack,
		Update:  nil,
		LastMes: -1,
	})
}

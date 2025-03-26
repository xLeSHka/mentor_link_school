package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
)

func MainMenuKeyboard(bot *Bot, telegram string) (tgbotapi.InlineKeyboardMarkup, error) {
	_, err := bot.UsersService.GetByTelegram(context.Background(), telegram)
	if err != nil && err.(*httpError.HTTPError).StatusCode == 500 {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	if err != nil && err.(*httpError.HTTPError).StatusCode == 404 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Регистрация", "Регистрация"),
			),
		), nil
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Авторизация", "Авторизация"),
		),
	), nil
}

func MainMenu(stack CallStack) CallStack {
	stack.Action = MainMenu
	data := userDatas[stack.ChatID]
	if data.User == nil {
		if stack.IsPrint {
			stack.IsPrint = false
			keyboard, err := MainMenuKeyboard(stack.Bot, stack.Data)
			if err != nil {
				log.Println(err)
				return stack
			}
			text := fmt.Sprintf("%s\n\n%s", MainMenuTemplate, MainMenuTextTemplate)
			if data.LastMes == -1 {

				msg := tgbotapi.NewMessage(stack.ChatID, text)
				msg.ReplyMarkup = keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return stack
				}
				data.LastMes = sended.MessageID
				log.Println("Set last msg to", sended.MessageID)
				log.Println(*userDatas[stack.ChatID])
				return stack
			} else {
				msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return stack
				}
				return stack
			}

		}
		if stack.Update != nil {
			if stack.Update.Message != nil {
				data.LastMes = -1
				stack.IsPrint = true
				return stack
			}
			if stack.Update.CallbackQuery != nil {
				switch stack.Update.CallbackQuery.Data {
				case "Регистрация":
					userDatas[stack.ChatID].User = &models.User{}
					return RegisterName(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						IsPrint: true,
						Parent:  &stack,
						Update:  nil,
						Data:    stack.Data,
					})
				case "Авторизация":
					userDatas[stack.ChatID].User = &models.User{}
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
	}
	data.LastMes = -1
	return AuthedMenu(CallStack{
		ChatID:  stack.ChatID,
		Bot:     stack.Bot,
		IsPrint: true,
		Parent:  &stack,
		Update:  nil,
		Data:    stack.Data,
	})
}

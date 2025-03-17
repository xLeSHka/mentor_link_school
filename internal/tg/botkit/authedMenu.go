package botkit

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"log"
	"net/http"
)

func AuthedMenuKeyboard(bot *Bot, telegram string) (tgbotapi.InlineKeyboardMarkup, error) {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать группу", "Создать группу"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мои группы", "Мои группы"),
		),
	), nil
}

func AuthedMenu(stack CallStack) CallStack {
	stack.Action = AuthedMenu
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.User.BIO == nil {
			bio := "Отсутствует"
			data.User.BIO = &bio
		}
		if data.User.AvatarURL != nil {
			err := avatar.GetUserAvatar(data.User, stack.Bot.MinioRepository)
			if err != nil {
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось получить ссылку на вашу аватарку!", ErrorMenuTemplate)))
				if err != nil {
					log.Println(err)
				}
				return stack
			}
			response, err := http.Get(*data.User.AvatarURL)
			if err != nil {
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось загрузить вашу аватарку!", ErrorMenuTemplate)))
				if err != nil {
					log.Println(err)
				}
				return stack
			}
			defer response.Body.Close()
			photoFileReader := tgbotapi.FileReader{
				Name:   "picture",
				Reader: response.Body,
			}
			_, err = stack.Bot.Api.Send(tgbotapi.NewPhoto(stack.ChatID, photoFileReader))
			if err != nil {
				log.Println(err)
				return stack
			}
		}
		text := fmt.Sprintf("%s\n\n%s", AuthedMenuTemplate, ProfileTextTemplate, data.User.ID, data.User.Name, data.User.BIO)
		msg := tgbotapi.NewMessage(stack.ChatID, text)

		keyboard, err := AuthedMenuKeyboard(stack.Bot, stack.Data)
		if err != nil {
			log.Println(err)
			return stack
		}
		msg.ReplyMarkup = keyboard
		_, err = stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return stack
	}

	return stack
}

package botkit

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
	"strings"
)

func EditUserKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Имя", "Имя"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("БИО", "БИО"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Аватар", "Аватар"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
}

func EditUser(stack CallStack) CallStack {
	stack.Action = EditUser
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.User.BIO == nil {
			bio := "Отсутствует"
			data.User.BIO = &bio
		}
		if stack.LastMes == -1 {
			if data.User.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.User.AvatarURL)
				if err != nil {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось получить ссылку на вашу аватарку!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err, 0)
					}
					return stack
				}
				avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
				response, err := http.Get(avatarURL)
				if err != nil {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось загрузить вашу аватарку!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err, -1)
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
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
			}
			removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
			msg := tgbotapi.NewMessage(stack.ChatID, "1")
			msg.ReplyMarkup = removeKeyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 1)
				return stack
			}
			resp, err := stack.Bot.Api.Request(tgbotapi.NewDeleteMessage(stack.ChatID, sended.MessageID))
			if err != nil || !resp.Ok {
				log.Println(err, 2)
				return stack
			}
			keyboard := EditUserKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO))
			msg = tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			stack.LastMes = sended.MessageID
		} else {
			keyboard := EditUserKeyboard()
			text := fmt.Sprintf("%s\n\n%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO), "Выберите что вы хотите отредактирвоать!")
			msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, text)
			msg.ReplyMarkup = &keyboard
			_, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
		}
		return stack
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "Имя":
				userDatas[stack.ChatID].Group = &models.Group{}
				return Chop(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: stack.LastMes,
					Data:    "Created1",
				})
			case "БИО":
				data.Size = 10
				data.Page = 0
				return Chop(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: stack.LastMes,
					Data:    "Created1",
				})
			case "Аватар":
				userDatas[stack.ChatID].Group = &models.Group{}
				return Chop(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: stack.LastMes,
					Data:    "Created1",
				})
			case "Выйти 🚪":
				userDatas[stack.ChatID].User = nil
				return ReturnOnParent(stack)
			default:
				return stack
			}
		}
	}
	return stack
}

package botkit

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
	"strings"
)

func AuthedMenuKeyboard(bot *Bot, telegram string) (tgbotapi.InlineKeyboardMarkup, error) {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать группу", "Создать группу"),
			tgbotapi.NewInlineKeyboardButtonData("Мои группы", "Мои группы"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Войти в группу", "Войти в группу"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать профиль", "Редактировать профиль"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выйти 🚪", "Выйти 🚪"),
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

		if data.LastMes == -1 {
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
				msg := tgbotapi.NewPhoto(stack.ChatID, photoFileReader)
				keyboard, err := AuthedMenuKeyboard(stack.Bot, stack.Data)
				if err != nil {
					log.Println(err, 3)
					return stack
				}
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\n%s", AuthedMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO))
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
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
			keyboard, err := AuthedMenuKeyboard(stack.Bot, stack.Data)
			if err != nil {
				log.Println(err, 3)
				return stack
			}
			text := fmt.Sprintf("%s\n\n%s", AuthedMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO))
			msg = tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			data.LastMes = sended.MessageID
		} else {
			if data.User.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.User.AvatarURL)
				if err != nil {
					data.LastMes = -1
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось получить ссылку на вашу аватарку!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err, 0)
					}
					return stack
				}
				avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
				response, err := http.Get(avatarURL)
				if err != nil {
					data.LastMes = -1
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
				baseMedia := tgbotapi.BaseInputMedia{
					Type:      "photo",
					Media:     photoFileReader,
					ParseMode: "markdown",
				}
				text := fmt.Sprintf("%s\n\n%s", AuthedMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO))
				baseMedia.Caption = text
				keyboard, err := AuthedMenuKeyboard(stack.Bot, stack.Data)
				if err != nil {
					log.Println(err, 3)
					return stack
				}

				msg := tgbotapi.EditMessageMediaConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:    stack.ChatID,
						MessageID: data.LastMes,
					},
					Media: tgbotapi.InputMediaPhoto{
						BaseInputMedia: baseMedia,
					},
				}
				msg.ReplyMarkup = &keyboard
				_, err = stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				return stack
			}
			keyboard, err := AuthedMenuKeyboard(stack.Bot, stack.Data)
			if err != nil {
				log.Println(err, 3)
				return stack
			}
			text := fmt.Sprintf("%s\n\n%s", AuthedMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO))
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
			msg.ReplyMarkup = &keyboard
			_, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
		}
		return stack
	}
	if stack.Update != nil {
		if stack.Update.Message != nil {
			data.LastMes = -1
			stack.IsPrint = true
			return stack
		}
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "Создать группу":
				userDatas[stack.ChatID].Group = &models.Group{}
				return CreateGroup(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Мои группы":
				data.Size = 10
				data.Page = 0
				return Groups(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Войти в группу":
				userDatas[stack.ChatID].Group = &models.Group{}
				return JoinToGroup(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Редактировать профиль":
				return EditUser(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Выйти 🚪":
				data.LastMes = -1
				userDatas[stack.ChatID].User = nil
				return ReturnOnParent(stack)
			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

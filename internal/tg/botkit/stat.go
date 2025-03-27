package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

func Stat(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Stat // Set self as current Action
	data := userDatas[stack.ChatID]
	stat, err := stack.Bot.GroupService.GetStat(context.Background(), data.Group.ID)
	if err != nil {
		_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
		if err != nil {
			log.Println(err)
			return ReturnOnParent(stack)
		}
		return ReturnOnParent(stack)
	}

	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
				if err != nil {
					data.LastMes = -1
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось получить ссылку на вашу аватарку!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err, 0)
					}
					return stack
				}

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
				text := fmt.Sprintf("%s\n\n%s", StatMenuTemplate, StatMenuTextTemplate(stat))
				keyboard := backInlineKeyboard()
				baseMedia.Caption = text
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
					log.Println(err, -21, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				return stack
			}
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s", StatMenuTemplate, StatMenuTextTemplate(stat)))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			_, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		} else {
			stack.IsPrint = false
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
				if err != nil {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНе удалось получить ссылку на вашу аватарку!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err, 0)
					}
					return stack
				}

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
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\n%s", StatMenuTemplate, StatMenuTextTemplate(stat))
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", StatMenuTemplate, StatMenuTextTemplate(stat)))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			data.LastMes = sended.MessageID
			return stack
		}
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					return ReturnOnParent(stack)
				}
			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

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
	"strings"
)

func JoinToGroup(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = JoinToGroup // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
			stack.IsPrint = false
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
				text := fmt.Sprintf("%s\n\nПожалуйста введите пригласительный код!", JoinMenuTemplate)
				baseMedia.Caption = text
				keyboard := backInlineKeyboard()

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
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\nПожалуйста введите пригласительный код!", JoinMenuTemplate))
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
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\nПожалуйста введите пригласительный код!", JoinMenuTemplate)
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
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
			data.LastMes = sended.MessageID
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
				data.LastMes = -1
				return stack
			}
		}
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			data.LastMes = -1
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
					data.LastMes = sended.MessageID
					return stack
				}
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate))
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
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
					_, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					return stack
				} else if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nВы уже состоите в этой организации!🚫", ErrorMenuTemplate))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					_, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else {
					msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate))
					keyboard := backInlineKeyboard()
					msg.ReplyMarkup = &keyboard
					_, err := stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				}
				return ReturnOnParent(stack)
			}
			data.LastMes = -1
			return Group(CallStack{
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

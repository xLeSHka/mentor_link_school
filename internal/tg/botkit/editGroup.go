package botkit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bachvtuan/mime2extension"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
)

func EditGroupKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Имя", "Имя"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Инвайт код", "Инвайт код"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Аватар", "Аватар"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
}

func EditGroup(stack CallStack) CallStack {
	stack.Action = EditGroup
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if stack.LastMes == -1 {
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
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

			keyboard := EditGroupKeyboard()
			text := fmt.Sprintf("%s\n\n%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode), "Выберите что вы хотите отредактирвоать!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			stack.LastMes = sended.MessageID
		} else {
			keyboard := EditGroupKeyboard()
			text := fmt.Sprintf("%s\n\n%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode), "Выберите что вы хотите отредактирвоать!")
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
				return EditGroupName(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: stack.LastMes,
					Data:    "Created1",
				})
			case "Инвайт код":
				code, err := stack.Bot.GroupService.UpdateInviteCode(context.Background(), data.Group.ID)
				if err != nil {
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nОрганизация не найдена!", ErrorMenuTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					} else {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					}
				}
				data.Group.InviteCode = &code
				keyboard := EditGroupKeyboard()
				text := fmt.Sprintf("%s\n\n%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode), "Выберите что вы хотите отредактирвоать!")
				msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err = stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 4)
					return stack
				}
				return stack
			case "Аватар":
				return EditGroupAvatar(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: stack.LastMes,
					Data:    "Created1",
				})
			case "⬅️ Назад":
				return ReturnOnParent(stack)
			default:
				return stack
			}
		}
	}
	return stack
}

func EditGroupName(stack CallStack) CallStack {
	stack.Action = EditGroupName
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if stack.LastMes == -1 {
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
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

			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое имя!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			stack.LastMes = sended.MessageID
		} else {
			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое имя!")
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
			case "⬅️ Назад":
				return ReturnOnParent(stack)
			}
		}
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			if len(msgText) > 120 {
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНевалидное имя!\nПожалуйста введите другое имя!", EditUserMenuTemplate))
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				stack.LastMes = sended.MessageID
				return stack
			}
			group, err := stack.Bot.GroupService.Edit(context.Background(), &models.Group{ID: data.Group.ID, Name: msgText})
			if err != nil {
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь не найден!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				}
			}
			data.Group = group
			return ReturnOnParent(stack)
		}
	}
	return stack
}

func EditGroupAvatar(stack CallStack) CallStack {
	stack.Action = EditGroupAvatar
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if stack.LastMes == -1 {
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
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

			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			stack.LastMes = sended.MessageID
		} else {
			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
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
			case "⬅️ Назад":
				return ReturnOnParent(stack)
			}
		}
		if stack.Update.Message.Document != nil {
			photo := stack.Update.Message.Document
			toSave, err := stack.Bot.Api.GetFile(tgbotapi.FileConfig{
				FileID: photo.FileID,
			})
			if err != nil {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось получить новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", stack.Bot.Api.Token, toSave.FilePath)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось подготовить запрос!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось загрузить новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			defer resp.Body.Close()
			buff := &bytes.Buffer{}
			_, err = io.Copy(buff, resp.Body)
			if err != nil {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось прочитать новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			mime := http.DetectContentType(buff.Bytes())
			_, ext := mime2extension.Extension(mime)
			f := &models.File{
				Filename: data.Group.ID.String() + "." + ext,
				File:     buff,
				Size:     int64(toSave.FileSize),
				Mimetype: mime,
			}
			_, cErr := stack.Bot.GroupService.UploadImage(context.Background(), f, data.Group.ID)
			if cErr != nil {
				log.Println(cErr)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось сохранить новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			data.Group.AvatarURL = &f.Filename
			return ReturnOnParent(stack)
		}
	}
	return stack
}

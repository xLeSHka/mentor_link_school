package botkit

import (
	"bytes"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"image"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func CreateGroup(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = CreateGroup // Set self as current Action
	data := userDatas[stack.ChatID]
	data.Group.ID = uuid.New()
	if stack.IsPrint {
		stack.IsPrint = false
		if data.User.AvatarURL != nil {
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя организации!", CreateGroupMenuTemplate, CreateGroupTextTemplate(data.Group.ID, "____", "____")))
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
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя организации!", CreateGroupMenuTemplate, CreateGroupTextTemplate(data.Group.ID, "____", "____")))
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
				sended, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНевалидне имя!\nМинимальная длина: %s\nМаксимальная длина: %s\nПожалуйста введите другое имя!", CreateGroupMenuTemplate, minL, maxL)))
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				data.LastMes = sended.MessageID
				return stack
			}
			data.Group.Name = msgText
			data.LastMes = -1
			return CreateGroupAvatar(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Update:  nil,
			})

		}
	}
	return stack
}
func CreateGroupAvatar(stack CallStack) CallStack {
	stack.Action = CreateGroupAvatar
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes == -1 {
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
				msg := tgbotapi.NewPhoto(stack.ChatID, photoFileReader)
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}

			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", CreateGroupMenuTemplate, "Загрузите новую аватарку без сжатия!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			data.LastMes = sended.MessageID
		} else {
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
				keyboard := backInlineKeyboard()
				text := fmt.Sprintf("%s\n\n%s", CreateGroupMenuTemplate, "Загрузите новую аватарку без сжатия!")
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
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				return stack
			}
			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", CreateGroupMenuTemplate, "Загрузите новую аватарку без сжатия!")
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
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
			data.LastMes = -1
			photo := stack.Update.Message.Document
			toSave, err := stack.Bot.Api.GetFile(tgbotapi.FileConfig{
				FileID: photo.FileID,
			})
			if err != nil {
				data.LastMes = -1
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
				data.LastMes = -1
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
				data.LastMes = -1
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
				data.LastMes = -1
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось прочитать новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			decodeBuff := bytes.NewBuffer(buff.Bytes())
			imgCfg, _, err := image.DecodeConfig(decodeBuff)
			if err != nil {
				data.LastMes = -1
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось декодировать изображение!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return stack
			}
			if len(buff.Bytes()) > 10<<20 || imgCfg.Height+imgCfg.Width > 10000 || imgCfg.Height/imgCfg.Width > 20 || imgCfg.Width/imgCfg.Height > 20 {
				data.LastMes = -1
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\n%s", ErrorMenuTemplate, "Изображение не подходит по формату!", checkPhoto(buff, imgCfg))))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return stack
			}
			ext := filepath.Ext(photo.FileName)
			mime := mime.TypeByExtension(ext)
			log.Println(ext, mime)
			f := &models.File{
				Filename: data.Group.ID.String() + ext,
				File:     buff,
				Size:     int64(toSave.FileSize),
				Mimetype: mime,
			}
			inviteCode, err := stack.Bot.GroupService.Create(context.Background(), data.Group, data.User.ID)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				data.LastMes = -1
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			data.Group.InviteCode = &inviteCode
			data.LastMes = -1
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
			return Group(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Update:  nil,
				Data:    "Created3",
			})
		}
	}
	return stack
}

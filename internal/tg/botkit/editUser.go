package botkit

import (
	"bytes"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
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
				keyboard := EditUserKeyboard()
				text := fmt.Sprintf("%s\n\n%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO), "Выберите что вы хотите отредактирвоать!")
				msg.Caption = text
				msg.ReplyMarkup = keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}

			keyboard := EditUserKeyboard()
			text := fmt.Sprintf("%s\n\n%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO), "Выберите что вы хотите отредактирвоать!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
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
				keyboard := EditUserKeyboard()
				text := fmt.Sprintf("%s\n\n%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO), "Выберите что вы хотите отредактирвоать!")

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
			keyboard := EditUserKeyboard()
			text := fmt.Sprintf("%s\n\n%s\n\n%s", EditUserMenuTemplate, ProfileTextTemplate(data.User.ID, data.User.Name, *data.User.BIO), "Выберите что вы хотите отредактирвоать!")
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
		if stack.Update.Message != nil {
			data.LastMes = -1
			stack.IsPrint = true
			return stack
		}
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "Имя":
				return EditName(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "БИО":
				return EditBIO(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Аватар":
				return EditAvatar(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "⬅️ Назад":
				return ReturnOnParent(stack)
			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

func EditName(stack CallStack) CallStack {
	stack.Action = EditName
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
				text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое имя!")
				msg.Caption = text
				msg.ReplyMarkup = keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
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
			data.LastMes = sended.MessageID
		} else {
			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое имя!")
			if data.User.AvatarURL != nil {
				msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 4)
					return stack
				}
			} else {

				msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 4)
					return stack
				}
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
			data.LastMes = -1
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
				data.LastMes = sended.MessageID
				return stack
			}
			user, err := stack.Bot.UsersService.Edit(context.Background(), data.User.ID, &models.User{Name: msgText})
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
			data.User = user
			return ReturnOnParent(stack)
		}
	}
	return stack
}

func EditBIO(stack CallStack) CallStack {
	stack.Action = EditBIO
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
				text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое БИО!")
				msg.Caption = text
				msg.ReplyMarkup = keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое БИО!")
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err, 4)
				return stack
			}
			data.LastMes = sended.MessageID
		} else {

			keyboard := backInlineKeyboard()
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Введите новое БИО!")
			if data.User.AvatarURL != nil {

				msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 4)
					return stack
				}
			} else {
				msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
				msg.ReplyMarkup = &keyboard
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 4)
					return stack
				}
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
			data.LastMes = -1
			msgText := stack.Update.Message.Text
			if len(msgText) > 120 {
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНевалидное БИО!\nПожалуйста введите другое БИО!", EditUserMenuTemplate))
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				data.LastMes = sended.MessageID
				return stack
			}
			user, err := stack.Bot.UsersService.Edit(context.Background(), data.User.ID, &models.User{BIO: &msgText})
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
			data.User = user
			return ReturnOnParent(stack)
		}
	}
	return stack
}
func checkPhoto(buff *bytes.Buffer, imgCfg image.Config) string {
	if len(buff.Bytes()) > 10<<20 {
		return "Размер изображения не может превышать 10 мегабайт!"
	}
	if imgCfg.Height+imgCfg.Width > 10000 {
		return "Сумма высоты и ширины изображения не может превышать 10000!"
	}
	if imgCfg.Height/imgCfg.Width > 20 || imgCfg.Width/imgCfg.Height > 20 {
		return "Соотношение сторон изображения не может быть больше 20!"
	}
	return "Неопознанная ошибка!"
}
func EditAvatar(stack CallStack) CallStack {
	stack.Action = EditAvatar
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
				text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
				msg.Caption = text
				msg.ReplyMarkup = keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
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
				text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
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
			text := fmt.Sprintf("%s\n\n%s", EditUserMenuTemplate, "Загрузите новую аватарку без сжатия!")
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
			decodeBuff := bytes.NewBuffer(buff.Bytes())
			imgCfg, _, err := image.DecodeConfig(decodeBuff)
			if err != nil {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось декодировать изображение!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return stack
			}
			if len(buff.Bytes()) > 10<<20 || imgCfg.Height+imgCfg.Width > 10000 || imgCfg.Height/imgCfg.Width > 20 || imgCfg.Width/imgCfg.Height > 20 {
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\n%s", ErrorMenuTemplate, "Изображение не подходит по формату!", checkPhoto(buff, imgCfg))))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return stack
			}
			ext := filepath.Ext(photo.FileName)
			mime := mime.TypeByExtension(ext)
			log.Println(ext, mime, buff.Len(), toSave.FileSize)
			f := &models.File{
				Filename: data.User.ID.String() + ext,
				File:     buff,
				Size:     int64(toSave.FileSize),
				Mimetype: mime,
			}
			_, cErr := stack.Bot.UsersService.UploadImage(context.Background(), f, data.User.ID)
			if cErr != nil {
				log.Println(cErr)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, "Не удалось сохранить новую аватарку!")))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				return ReturnOnParent(stack)
			}
			data.User.AvatarURL = &f.Filename
			return ReturnOnParent(stack)
		}
	}
	return stack
}

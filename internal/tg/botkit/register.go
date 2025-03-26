package botkit

import (
	"bytes"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/password"
	"image"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"unicode"
)

func backInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
}
func RegisterName(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = RegisterName     // Set self as current Action
	data := userDatas[stack.ChatID] // User data
	data.User.ID = uuid.New()
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes == -1 {
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, "____", "____", "____")))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			data.LastMes = sended.MessageID
			return stack
		} else {
			log.Println(data.LastMes)
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, "____", "____", "____")))
			keyboard := backInlineKeyboard()
			msg.ReplyMarkup = &keyboard
			_, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		}
	} else {
		// Processing a message
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
			}
		}
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			data.LastMes = -1
			if len(msgText) > 120 {
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nНевалидное имя!\nПожалуйста введите другое имя!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, "____", "____", "____")))
				msg.ReplyMarkup = backInlineKeyboard()
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				return stack
			}
			data.User.Name = msgText
			data.User.Telegram = stack.Update.Message.From.UserName
			return RegisterPassword(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Data:    stack.Data,
			})

		}

	}

	// Repeat self
	return stack
}

func RegisterPassword(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = RegisterPassword // Set self as current Action
	data := userDatas[stack.ChatID] // User data
	if stack.IsPrint {
		stack.IsPrint = false
		// Print UI
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите пароль!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, data.User.Name, data.User.Telegram, "____")))

		msg.ReplyMarkup = backInlineKeyboard()
		_, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			return ReturnOnParent(stack)
		}
		// Remove previous Keyboard or set self
		return stack
	}
	if stack.Update != nil {
		// Processing a message
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					return ReturnOnParent(stack)
				}
			}
		}
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			data.LastMes = -1
			hasLower, hasUpper, hasDigit, hasSymbol := validatePassword(msgText)
			if len(msgText) > 60 || len(msgText) < 8 || !(hasLower && hasUpper && hasDigit && hasSymbol) {
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nНевалидный пароль!\nПожалуйста введите другой пароль!", RegisterMenuTemplate, ValidatePasswordTemplate(len(msgText) >= 8, len(msgText) <= 60, hasLower, hasUpper, hasDigit, hasSymbol)))
				msg.ReplyMarkup = backInlineKeyboard()
				_, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				return stack
			}
			encoded, err := password.Encrypt([]byte(msgText), stack.Bot.CryptoKey)
			if err != nil {
				log.Println(err)
				_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				return stack
			}
			data.User.Password = encoded
			data.User.TelegramID = &stack.ChatID

			return RegisterAvatar(CallStack{
				ChatID:  stack.ChatID,
				Bot:     stack.Bot,
				IsPrint: true,
				Parent:  &stack,
				Data:    stack.Data,
			})

		}
	}
	return stack
}
func validatePassword(password string) (bool, bool, bool, bool) {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.In(char, unicode.Symbol, unicode.Punct):
			hasSpecial = true
		}
	}

	return hasLower, hasUpper, hasDigit, hasSpecial
}

func RegisterAvatar(stack CallStack) CallStack {
	stack.Action = RegisterAvatar
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false

		keyboard := backInlineKeyboard()
		text := fmt.Sprintf("%s\n\n%s\n\nПожалуйста загрузите аватарку без сжатия!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, data.User.Name, data.User.Telegram, "****"))
		msg := tgbotapi.NewMessage(stack.ChatID, text)
		msg.ReplyMarkup = keyboard
		sended, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err, 4)
			return stack
		}
		data.LastMes = sended.MessageID
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
			_, err = stack.Bot.UsersService.Register(context.Background(), data.User)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nНельзя регистрировать два аккаунта с одним телеграммом!🚫", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						ReturnOnParent(*stack.Parent)
					}
				} else {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						ReturnOnParent(*stack.Parent)
					}
				}
				return ReturnOnParent(*stack.Parent)
			}
			user, err := stack.Bot.UsersService.GetByTelegram(context.Background(), stack.Data)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь с вашим телеграм не найден!🤨🔎", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						ReturnOnParent(*stack.Parent)
					}
				} else {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						ReturnOnParent(*stack.Parent)
					}
				}
				return ReturnOnParent(stack)
			}
			userDatas[stack.ChatID].User = user
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
			return ReturnOnParent(*stack.Parent.Parent)
		}
	}
	return stack
}

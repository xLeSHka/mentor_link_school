package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
)

func LoginPassword(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = LoginPassword // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		log.Println(data.LastMes)
		stack.IsPrint = false
		if data.LastMes == -1 {
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПожалуйста введите пароль!", LoginMenuTemplate))
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
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\nПожалуйста введите пароль!", LoginMenuTemplate))
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
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					log.Println("Back", data.LastMes)
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
			}
		}
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			data.LastMes = -1
			hasLower, hasUpper, hasDigit, hasSymbol := validatePassword(msgText)
			if len(msgText) > 60 || len(msgText) < 8 || !(hasLower && hasUpper && hasDigit && hasSymbol) {
				msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nНевалидный пароль!\nПожалуйста введите другой пароль!", LoginMenuTemplate, ValidatePasswordTemplate(len(msgText) >= 8, len(msgText) <= 60, hasLower, hasUpper, hasDigit, hasSymbol)))
				keyboard := backInlineKeyboard()
				msg.ReplyMarkup = &keyboard
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				stack.Update = nil
				data.LastMes = sended.MessageID
				return stack
			}
			_, err := stack.Bot.UsersService.Login(context.Background(), stack.Data, msgText)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				data.LastMes = -1
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь с вашим телеграм не найден!🤨🔎", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else if err.(*httpError.HTTPError).StatusCode == http.StatusUnauthorized {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПароль введен не верно!🚫\nПожалуйста введите пароль!", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return stack
					}
				} else {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				}
				return ReturnOnParent(stack)
			}
			user, err := stack.Bot.UsersService.GetByTelegram(context.Background(), stack.Data)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				data.LastMes = -1
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь с вашим телеграм не найден!🤨🔎", ErrorMenuTemplate)))
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
				return ReturnOnParent(stack)
			}
			_, err = stack.Bot.UsersService.Edit(context.Background(), user.ID, &models.User{TelegramID: &stack.ChatID})
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].User = nil
				data.LastMes = -1
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь с вашим телеграм не найден!🤨🔎", ErrorMenuTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
				} else if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь с таким ID уже зарегистрирован! 🚫", ErrorMenuTemplate)))
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
				return ReturnOnParent(stack)
			}
			user.TelegramID = &stack.ChatID
			userDatas[stack.ChatID].User = user
			return ReturnOnParent(stack)

		}
	}
	return stack
}

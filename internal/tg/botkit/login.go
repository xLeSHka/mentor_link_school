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
	if stack.IsPrint {
		stack.IsPrint = false
		// Print UI
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПожалуйста введите пароль!", LoginMenuTemplate))

		msg.ReplyMarkup = backButton()
		_, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			userDatas[stack.ChatID].User = nil
			return ReturnOnParent(stack)
		}
		// Remove previous Keyboard or set self
		return stack
	}
	if stack.Update != nil {
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			switch msgText {
			case "⬅️ Назад":
				{
					userDatas[stack.ChatID].User = nil
					return ReturnOnParent(stack)
				}
			default:
				hasLower, hasUpper, hasDigit, hasSymbol := validatePassword(msgText)
				if len(msgText) > 60 || len(msgText) < 8 || !(hasLower && hasUpper && hasDigit && hasSymbol) {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nНевалидный пароль!\nПожалуйста введите другой пароль!", RegisterMenuTemplate, ValidatePasswordTemplate(len(msgText) >= 8, len(msgText) <= 60, hasLower, hasUpper, hasDigit, hasSymbol))))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.Update = nil
					return stack
				}
				_, err := stack.Bot.UsersService.Login(context.Background(), stack.Data, msgText)
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].User = nil
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
	}
	return stack
}

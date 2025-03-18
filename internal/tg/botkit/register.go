package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/password"
	"log"
	"net/http"
	"unicode"
)

func backButton() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Назад"),
		),
	)
}

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
		// Print UI
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nПожалуйста введите имя!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, "____", "____", "____")))

		msg.ReplyMarkup = backButton()
		_, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			userDatas[stack.ChatID].User = nil
			return ReturnOnParent(stack)
		}
		// Remove previous Keyboard or set self
		return stack
	} else {
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			switch msgText {
			case "⬅️ Назад":
				userDatas[stack.ChatID].User = nil
				return ReturnOnParent(stack)
			default:
				if len(msgText) > 120 {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s\n\nНевалидное имя!\nПожалуйста введите другое имя!", RegisterMenuTemplate, RegisterMenuTextTemplate(data.User.ID, "____", "____", "____"))))
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

		msg.ReplyMarkup = backButton()
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
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			switch msgText {
			case "⬅️ Назад":
				{
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
				return ReturnOnParent(*stack.Parent)
			}
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

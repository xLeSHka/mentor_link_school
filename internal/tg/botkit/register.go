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

var backButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⬅️ Back"),
	),
)

func RegisterName(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = RegisterName     // Set self as current Action
	data := userDatas[stack.ChatID] // User data
	data.User.ID = uuid.New()
	data.User.Telegram = stack.Update.Message.From.UserName
	if stack.IsPrint {
		stack.IsPrint = false
		// Print UI
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Регистрационное меню]\n\nВы 🫵\nID: 🆔%s\nИмя: %s 🪪\nТелеграм: %s \nПароль: %s 🔑\n\nПожалуйста введите имя!", data.User.ID, "____", data.User.Telegram, "____"))

		msg.ReplyMarkup = backButton
		_, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			return ReturnOnParent(stack)
		}
		// Remove previous Keyboard or set self
		return stack
	} else {
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			switch msgText {
			case "⬅️ Back":
				return ReturnOnParent(stack)
			default:
				if len(msgText) > 120 {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Регистрационное меню]\n\nВы 🫵\nID: 🆔%s\nИмя: %s 🪪\nТелеграм: %s \nПароль: %s 🔑\n\nНевалидное имя!\nПожалуйста введите другое имя!", data.User.ID, "____", data.User.Telegram, "____")))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.Update = nil
					return RegisterName(stack)
				}
				data.User.Name = msgText
				return RegisterPassword(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
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
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Регистрационное меню]\n\nВы 🫵\nID: 🆔%s\nИмя: %s 🪪\nТелеграм: %s \nПароль: %s 🔑\n\nПожалуйста введите пароль!", data.User.ID, data.User.Name, data.User.Telegram, "____"))

		msg.ReplyMarkup = backButton
		_, err := stack.Bot.Api.Send(msg)
		if err != nil {
			log.Println(err)
			return ReturnOnParent(stack)
		}
		// Remove previous Keyboard or set self
		return stack
	} else {
		// Processing a message
		if stack.Update.Message != nil {
			msgText := stack.Update.Message.Text
			switch msgText {
			case "back":
				{
					return ReturnOnParent(stack)
				}
			default:
				if len(msgText) > 60 || len(msgText) < 8 || !validatePassword(msgText) {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Регистрационное меню]\n\nВы 🫵\nID: 🆔%s\nИмя: %s 🪪\nТелеграм: %s \nПароль: %s 🔑\n\nНевалидный пароль!\nПожалуйста введите другой пароль!", data.User.ID, data.User.Name, data.User.Telegram, "____")))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					stack.Update = nil
					return RegisterPassword(stack)
				}
				encoded, err := password.Encrypt([]byte(msgText), stack.Bot.CryptoKey)
				if err != nil {
					log.Println(err)
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Ошибка]\n\nПриносим свои извинения, произошла непредвиденная ошибка! 🥺🙏")))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}
					return RegisterPassword(stack)
				}
				data.User.Password = encoded
				_, err = stack.Bot.UsersService.Register(context.Background(), data.User)
				if err != nil {
					log.Println(err)
					if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Ошибка]\n\nНельзя регистрировать два аккаунта с одним телеграммом! 🚫")))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					} else {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[Ошибка]\n\nПриносим свои извинения, произошла непредвиденная ошибка! 🥺🙏")))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					}
					return MainMenu(CallStack{
						ChatID:  stack.ChatID,
						Bot:     stack.Bot,
						IsPrint: true,
					})
				}
				return Chop(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
				})
			}
		}
	}
	return stack
}
func validatePassword(password string) bool {
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

	return hasLower && hasUpper && hasDigit && hasSpecial
}

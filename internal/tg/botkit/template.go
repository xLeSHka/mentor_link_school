package botkit

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunTemplate(stack CallStack) CallStack {
	return Chop(stack)          // delete or comment out after finishing work
	stack.Action = RunTemplate  // Set self as current Action
	_ = userDatas[stack.ChatID] // User data

	if stack.IsPrint {
		stack.IsPrint = false // delete or comment out if print repeated
		// Print UI
		msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("[menu]"))
		mainMenuInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "back"),
			),
		)
		msg.ReplyMarkup = mainMenuInlineKeyboard
		_, _ = stack.Bot.Api.Send(msg)
		// Remove previous Keyboard or set self
		mainMenuKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg = tgbotapi.NewMessage(stack.ChatID, "")
		msg.ReplyMarkup = mainMenuKeyboard
		_, _ = stack.Bot.Api.Send(msg)
		return stack
	} else {
		// Processing a message
		if stack.Update.Message != nil {
			switch stack.Update.Message.Text {
			case "back":
				{
					return ReturnOnParent(stack)
				}
			}
		}
		if stack.Update != nil {
			// Processing a message
			if stack.Update.CallbackQuery != nil {
				switch stack.Update.CallbackQuery.Data {
				case "⬅️ Back":
					{
						return ReturnOnParent(stack)
					}
				}
			}
		}
		if stack.Update.Message.IsCommand() {
			switch stack.Update.Message.Command() {
			case "back":
				{
					return ReturnOnParent(stack)
				}
			}
		}
	}

	// Repeat self
	return stack
}

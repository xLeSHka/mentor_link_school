package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"strings"
)

func mapMember(role *models.User) (string, string) {
	builder := strings.Builder{}

	builder.WriteString(role.Name)
	for _, role := range role.Roles {
		switch role.Role {
		case "owner":
			builder.WriteString("🧑‍💼")
		case "mentor":
			builder.WriteString("🧑‍🏫")
		case "student":
			builder.WriteString("👨‍🎓")
		}
	}
	return builder.String(), role.ID.String()
}
func MembersKeyboard(bot *Bot, groupID uuid.UUID) (tgbotapi.InlineKeyboardMarkup, error) {
	members, err := bot.GroupService.GetMembers(context.Background(), groupID)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(members))
	for _, member := range members {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(mapMember(member)),
			),
		)
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
func Members(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Members // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		if stack.LastMes != -1 {

			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MembersMenuTemplate, MembersMenuTextTemplate()))
			keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID)
			if err != nil {
				data.Profile = nil
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}

				return ReturnOnParent(stack)
			}
			msg.ReplyMarkup = &keyboard
			_, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				data.Profile = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		} else {
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", MembersMenuTemplate, MembersMenuTextTemplate()))
			keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID)
			if err != nil {
				data.Profile = nil
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}

				return ReturnOnParent(stack)
			}
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				data.Profile = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			stack.LastMes = sended.MessageID
			stack.Parent.LastMes = sended.MessageID
			return stack
		}
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					data.Profile = nil
					return ReturnOnParent(stack)
				}
			default:
				id := uuid.MustParse(msgText)
				data.Profile = &models.User{ID: id}
				return Profile(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
					LastMes: -1,
					Data:    "Created1",
				})
			}
		}
	}
	return stack
}

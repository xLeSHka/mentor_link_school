package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
	"strings"
)

func mapGroup(group *models.GroupWithRoles) (string, string) {
	builder := strings.Builder{}

	builder.WriteString(group.Group.Name)
	for _, role := range group.MyRoles {
		switch role {
		case "owner":
			builder.WriteString("🧑‍💼")
		case "mentor":
			builder.WriteString("🧑‍🏫")
		case "student":
			builder.WriteString("👨‍🎓")
		}
	}
	return builder.String(), group.GroupID.String()
}
func GroupsKeyboard(bot *Bot, userID uuid.UUID, page, size int) (tgbotapi.InlineKeyboardMarkup, error) {
	groups, total, err := bot.UsersService.GetGroups(context.Background(), userID, page, size)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(groups))
	for _, member := range groups {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(mapGroup(member)),
			),
		)
	}
	log.Println(total, page, size)
	navigation := tgbotapi.NewInlineKeyboardRow()
	if page > 0 {
		navigation = append(navigation, tgbotapi.NewInlineKeyboardButtonData("⬅️", "Влево"))
	}
	if total-int64(page*size+size) > 0 {
		navigation = append(navigation, tgbotapi.NewInlineKeyboardButtonData("➡️", "Вправо"))
	}
	if len(navigation) > 0 {
		rows = append(rows, navigation)
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
func Groups(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Groups // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
			text := fmt.Sprintf("%s\n\nПожалуйста выберите организацию!", GroupsMenuTemplate)
			stack.IsPrint = false
			keyboard, err := GroupsKeyboard(stack.Bot, data.User.ID, data.Page, data.Size)
			if err != nil {
				data.LastMes = -1
				data.Group = nil
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}

				return ReturnOnParent(stack)
			}
			if data.User.AvatarURL != nil {
				msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)

				msg.ReplyMarkup = &keyboard
				_, err = stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 1)
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
				// Remove previous Keyboard or set self
				return stack
			} else {

				// Print UI
				msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)

				msg.ReplyMarkup = &keyboard
				_, err = stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, 1)
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
				// Remove previous Keyboard or set self
				return stack
			}
		} else {
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПожалуйста выберите организацию!", GroupsMenuTemplate))
			keyboard, err := GroupsKeyboard(stack.Bot, data.User.ID, data.Page, data.Size)
			if err != nil {
				data.Group = nil
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
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			data.LastMes = sended.MessageID
			// Remove previous Keyboard or set self
			return stack
		}
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "Влево":
				{
					data.Page -= 1
					keyboard, err := GroupsKeyboard(stack.Bot, data.User.ID, data.Page, data.Size)
					if err != nil {
						data.Group = nil
						data.LastMes = -1
						_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
						if err != nil {
							log.Println(err)

							return ReturnOnParent(stack)
						}

						return ReturnOnParent(stack)
					}
					msg := tgbotapi.NewEditMessageReplyMarkup(stack.ChatID, data.LastMes, keyboard)
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Group = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "Вправо":
				{
					data.Page += 1
					keyboard, err := GroupsKeyboard(stack.Bot, data.User.ID, data.Page, data.Size)
					if err != nil {
						data.LastMes = -1
						data.Group = nil
						_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}

						return ReturnOnParent(stack)
					}
					msg := tgbotapi.NewEditMessageReplyMarkup(stack.ChatID, data.LastMes, keyboard)
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Group = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "⬅️ Назад":
				{
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
			default:
				groupID := uuid.MustParse(msgText)
				group, err := stack.Bot.GroupService.GetGroupByID(context.Background(), groupID)
				if err != nil {
					log.Println(err)
					userDatas[stack.ChatID].Group = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {

						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nОрганизация не найдена!🤨🔎", ErrorMenuTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
						return stack
					} else {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					}
					return ReturnOnParent(stack)
				}
				roles, _ := stack.Bot.GroupService.GetRoles(context.Background(), data.User.ID, groupID)
				isOwner := false
				for _, role := range roles {
					if role.Role == "owner" {
						isOwner = true
					}
				}
				if !isOwner {
					group.InviteCode = nil
				}
				userDatas[stack.ChatID].Group = group
				return Group(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			}
		}
	}
	return stack
}

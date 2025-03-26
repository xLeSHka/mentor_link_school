package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
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
func MembersKeyboard(bot *Bot, groupID uuid.UUID, page, size int) (tgbotapi.InlineKeyboardMarkup, error) {
	members, total, err := bot.GroupService.GetMembers(context.Background(), groupID, page, size)
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
func Members(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Members // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
			keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID, data.Page, data.Size)
			if err != nil {
				data.LastMes = -1
				data.Profile = nil
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}

				return ReturnOnParent(stack)
			}
			stack.IsPrint = false
			text := fmt.Sprintf("%s\n\n%s", MembersMenuTemplate, MembersMenuTextTemplate())
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
					log.Println(err, -21, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				return stack
			}

			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)

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
			text := fmt.Sprintf("%s\n\n%s", MembersMenuTemplate, MembersMenuTextTemplate())
			stack.IsPrint = false
			// Print UI
			keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID, data.Page, data.Size)
			if err != nil {
				data.Profile = nil
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}

				return ReturnOnParent(stack)
			}
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

				msg.ReplyMarkup = &keyboard
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			msg := tgbotapi.NewMessage(stack.ChatID, text)
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				data.Profile = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			data.LastMes = sended.MessageID
			return stack
		}
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
			case "Влево":
				{
					data.Page -= 1
					keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID, data.Page, data.Size)
					if err != nil {
						data.Profile = nil
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
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "Вправо":
				{
					data.Page += 1
					keyboard, err := MembersKeyboard(stack.Bot, data.Group.ID, data.Page, data.Size)
					if err != nil {
						data.Profile = nil
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
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
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
				})
			}
		}
	}
	return stack
}

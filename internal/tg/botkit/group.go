package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"log"
	"net/http"
	"strings"
)

func GroupKeyboard(bot *Bot, userID, groupID uuid.UUID) (tgbotapi.InlineKeyboardMarkup, error) {
	roles, err := bot.GroupService.GetRoles(context.Background(), userID, groupID)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	isStudent, isMentor, isOwner := false, false, false
	for _, role := range roles {
		switch role.Role {
		case "student":
			isStudent = true
		case "mentor":
			isMentor = true
		case "owner":
			isOwner = true
		}
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	if isOwner {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Члены организации 👥", "Члены организации 👥"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Редактировать организацию", "Редактировать организацию"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Статистика организации 📈", "Статистика организации 📈"),
			),
		)
	}
	if isMentor {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Мои студенты 👨‍🎓", "Мои студенты 👨‍🎓"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Входящие заявки 📩", "Входящие заявки 📩"),
			),
		)
	}
	if isStudent {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Мои менторы 🧑‍🏫", "Мои менторы 🧑‍🏫"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Исходящие заявки 📤", "Исходящие заявки 📤"),
				tgbotapi.NewInlineKeyboardButtonData("Доступные менторы", "Доступные менторы"),
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
func Group(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Group // Set self as current Action
	data := userDatas[stack.ChatID]
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
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
				text := fmt.Sprintf("%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode))
				keyboard, err := GroupKeyboard(stack.Bot, data.User.ID, data.Group.ID)
				if err != nil {
					data.LastMes = -1
					userDatas[stack.ChatID].Group = nil
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой группы не существует или вы в ней не состоите!", ErrorMenuTemplate)))
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
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode)))
			keyboard, err := GroupKeyboard(stack.Bot, data.User.ID, data.Group.ID)
			if err != nil {
				data.LastMes = -1
				userDatas[stack.ChatID].Group = nil
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой группы не существует или вы в ней не состоите!", ErrorMenuTemplate)))
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
			msg.ReplyMarkup = &keyboard
			_, err = stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		} else {
			stack.IsPrint = false
			if data.Group.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Group.AvatarURL)
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
				keyboard, err := GroupKeyboard(stack.Bot, data.User.ID, data.Group.ID)
				if err != nil {
					userDatas[stack.ChatID].Group = nil
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой группы не существует или вы в ней не состоите!", ErrorMenuTemplate)))
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
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode))
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", GroupMenuTemplate, GroupTextTemplate(data.Group.ID, data.Group.Name, data.Group.InviteCode)))
			keyboard, err := GroupKeyboard(stack.Bot, data.User.ID, data.Group.ID)
			if err != nil {
				userDatas[stack.ChatID].Group = nil
				if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
					_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой группы не существует или вы в ней не состоите!", ErrorMenuTemplate)))
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
			msg.ReplyMarkup = &keyboard
			sended, err := stack.Bot.Api.Send(msg)
			if err != nil {
				log.Println(err)
				userDatas[stack.ChatID].Group = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			data.LastMes = sended.MessageID
			return stack
		}
	}
	if stack.Update != nil {
		if stack.Update.CallbackQuery != nil {
			msgText := stack.Update.CallbackQuery.Data
			switch msgText {
			case "⬅️ Назад":
				{
					userDatas[stack.ChatID].Group = nil
					return ReturnOnParent(stack)
				}
			case "Члены организации 👥":
				data.Size = 10
				data.Page = 0
				return Members(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Редактировать организацию":
				return EditGroup(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Статистика организации 📈":
				return Stat(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Мои студенты 👨‍🎓":
				data.Size = 10
				data.Page = 0
				return Students(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Мои менторы 🧑‍🏫":
				data.Size = 10
				data.Page = 0
				return Mentors(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Доступные менторы":
				data.Size = 10
				data.Page = 0
				return AvailableMentors(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Исходящие заявки 📤":
				data.Size = 10
				data.Page = 0
				return SendedRequests(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			case "Входящие заявки 📩":
				data.Size = 10
				data.Page = 0
				return RecievedRequests(CallStack{
					ChatID:  stack.ChatID,
					Bot:     stack.Bot,
					IsPrint: true,
					Parent:  &stack,
					Update:  nil,
				})
			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

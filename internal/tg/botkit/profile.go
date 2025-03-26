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

func ProfileKeyboard(roles []*models.Role, userID uuid.UUID, bot *Bot) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	userRoles, err := bot.GroupService.GetRoles(context.Background(), userID, roles[0].GroupID)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	for _, userRole := range userRoles {
		switch userRole.Role {

		case "owner":

			mentor, student := "Добавить роль 🧑‍🏫", "Добавить роль 👨‍🎓"
			for _, role := range roles {
				switch role.Role {
				case "student":
					student = "Удалить роль 👨‍🎓"
				case "mentor":
					mentor = "Удалить роль 🧑‍🏫"
				}
			}
			rows = append(rows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(student, student),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(mentor, mentor),
				),
			)
		case "mentor":
		case "student":
		}
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
func Profile(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Profile // Set self as current Action
	data := userDatas[stack.ChatID]
	profile, err := stack.Bot.UsersService.GetByID(context.Background(), data.Profile.ID)

	if err != nil {
		data.LastMes = -1
		data.Profile = nil
		if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
			_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакой пользователь не существует!", ErrorMenuTemplate)))
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

	}
	roles, err := stack.Bot.GroupService.GetRoles(context.Background(), data.Profile.ID, data.Group.ID)
	if err != nil {
		data.LastMes = -1
		data.Profile = nil
		if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
			_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nУ этого пользователя нет ролей в этой организации!", ErrorMenuTemplate)))
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

	}
	data.Profile = profile
	if stack.IsPrint {
		stack.IsPrint = false
		if data.LastMes != -1 {
			if data.Profile.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Profile.AvatarURL)
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
				keyboard, err := ProfileKeyboard(roles, data.User.ID)
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
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))
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
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}

				return stack
			}
			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
			keyboard, err := ProfileKeyboard(roles, data.User.ID)
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
			if data.Profile.AvatarURL != nil {
				avatarURL, err := stack.Bot.MinioRepository.GetImage(*data.Profile.AvatarURL)
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
				keyboard, err := ProfileKeyboard(roles, data.User.ID)
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
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
			keyboard, err := ProfileKeyboard(roles, data.User.ID)
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
			case "⬅️ Назад":
				{
					data.Profile = nil
					return ReturnOnParent(stack)
				}
			case "Добавить роль 🧑‍🏫":
				err := stack.Bot.GroupService.AddRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "mentor",
				})
				if err != nil {
					data.LastMes = -1
					data.Profile = nil
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь не найден!", ErrorMenuTemplate)))
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

				}
				roles, _ = stack.Bot.GroupService.GetRoles(context.Background(), data.Profile.ID, data.Group.ID)
				keyboard, err := ProfileKeyboard(roles, data.User.ID)
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
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))
				if data.Profile.AvatarURL != nil {
					msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)

					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				} else {

					msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)

					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "Удалить роль 🧑‍🏫":
				err := stack.Bot.GroupService.RemoveRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "mentor",
				})
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь не найден!", ErrorMenuTemplate)))
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
				}
				roles, _ = stack.Bot.GroupService.GetRoles(context.Background(), data.Profile.ID, data.Group.ID)
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))

				keyboard, err := ProfileKeyboard(roles, data.User.ID)
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}

					return ReturnOnParent(stack)
				}
				if data.Profile.AvatarURL != nil {
					msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				} else {

					msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "Добавить роль 👨‍🎓":
				err := stack.Bot.GroupService.AddRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "student",
				})
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь не найден!", ErrorMenuTemplate)))
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

				}
				roles, _ = stack.Bot.GroupService.GetRoles(context.Background(), data.Profile.ID, data.Group.ID)

				keyboard, err := ProfileKeyboard(roles, data.User.ID)
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}

					return ReturnOnParent(stack)
				}
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))
				if data.Profile.AvatarURL != nil {
					msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				} else {

					msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}
			case "Удалить роль 👨‍🎓":
				err := stack.Bot.GroupService.RemoveRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "student",
				})
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nПользователь не найден!", ErrorMenuTemplate)))
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
				}
				roles, _ = stack.Bot.GroupService.GetRoles(context.Background(), data.Profile.ID, data.Group.ID)

				keyboard, err := ProfileKeyboard(roles, data.User.ID)
				if err != nil {
					data.Profile = nil
					data.LastMes = -1
					_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}

					return ReturnOnParent(stack)
				}
				text := fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles))
				if data.Profile.AvatarURL != nil {
					msg := tgbotapi.NewEditMessageCaption(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				} else {
					msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, text)
					msg.ReplyMarkup = &keyboard
					_, err = stack.Bot.Api.Send(msg)
					if err != nil {
						log.Println(err)
						data.Profile = nil
						return ReturnOnParent(stack)
					}
					return stack
				}

			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

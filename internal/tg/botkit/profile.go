package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"log"
	"net/http"
	"strings"
)

func ProfileKeyboard(roles []*models.Role) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)

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
		if stack.LastMes != -1 {

			stack.IsPrint = false
			// Print UI
			msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
			keyboard, err := ProfileKeyboard(roles)
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
				_, err = stack.Bot.Api.Send(tgbotapi.NewPhoto(stack.ChatID, photoFileReader))
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
			}
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
			keyboard, err := ProfileKeyboard(roles)
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
			case "Добавить роль 🧑‍🏫":
				err := stack.Bot.GroupService.AddRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "mentor",
				})
				if err != nil {
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
				msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
				keyboard, err := ProfileKeyboard(roles)
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
				return stack
			case "Удалить роль 🧑‍🏫":
				err := stack.Bot.GroupService.RemoveRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "mentor",
				})
				if err != nil {
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
				msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
				keyboard, err := ProfileKeyboard(roles)
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
				return stack
			case "Добавить роль 👨‍🎓":
				err := stack.Bot.GroupService.AddRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "student",
				})
				if err != nil {
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
				msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
				keyboard, err := ProfileKeyboard(roles)
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
				return stack
			case "Удалить роль 👨‍🎓":
				err := stack.Bot.GroupService.RemoveRole(context.Background(), &models.Role{
					UserID:  profile.ID,
					GroupID: data.Group.ID,
					Role:    "student",
				})
				if err != nil {
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
				msg := tgbotapi.NewEditMessageText(stack.ChatID, stack.LastMes, fmt.Sprintf("%s\n\n%s", MemberMenuTemplate, MemberMenuTextTemplate(profile.ID, profile.Name, profile.BIO, roles)))
				keyboard, err := ProfileKeyboard(roles)
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
				return stack
			default:
				return stack
			}
		}
	}
	return stack
}

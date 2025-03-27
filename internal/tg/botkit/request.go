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

func RequestKeyboard(role string) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	if role == "mentor" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✔️", "✔️"),
			tgbotapi.NewInlineKeyboardButtonData("❌", "❌"),
		))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "⬅️ Назад"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
func Request(stack CallStack) CallStack {
	//return Chop(stack)              // delete or comment out after finishing work
	stack.Action = Request // Set self as current Action
	data := userDatas[stack.ChatID]
	request, err := stack.Bot.StudentService.GetRequestByID(context.Background(), data.Req.ID, data.Group.ID)

	if err != nil {
		data.LastMes = -1
		data.Req = nil
		if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
			_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nТакого запроса не существует!", ErrorMenuTemplate)))
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
	data.Req = request
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
				keyboard, err := RequestKeyboard(stack.Data)
				if err != nil {
					data.LastMes = -1
					data.Req = nil
					_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}

					return ReturnOnParent(stack)
				}
				text := fmt.Sprintf("%s\n\n%s", RequestMenuTemplate, RequestMenuTextTemplate(data.Req))
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
			msg := tgbotapi.NewEditMessageText(stack.ChatID, data.LastMes, fmt.Sprintf("%s\n\n%s", RequestMenuTemplate, RequestMenuTextTemplate(data.Req)))
			keyboard, err := RequestKeyboard(stack.Data)
			if err != nil {
				data.LastMes = -1
				data.Req = nil
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
				data.Req = nil
				return ReturnOnParent(stack)
			}
			// Remove previous Keyboard or set self
			return stack
		} else {
			stack.IsPrint = false
			// Print UI
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
				keyboard, err := RequestKeyboard(stack.Data)
				if err != nil {
					data.Req = nil
					_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", ErrorMenuTemplate, InternalErrorTextTemplate)))
					if err != nil {
						log.Println(err)
						return ReturnOnParent(stack)
					}

					return ReturnOnParent(stack)
				}
				msg.ReplyMarkup = &keyboard
				text := fmt.Sprintf("%s\n\n%s", RequestMenuTemplate, RequestMenuTextTemplate(data.Req))
				msg.Caption = text
				sended, err := stack.Bot.Api.Send(msg)
				if err != nil {
					log.Println(err, -2, "ChatID ", stack.ChatID, "FileReader ", photoFileReader, "Url", avatarURL)
					return stack
				}
				data.LastMes = sended.MessageID
				return stack
			}
			msg := tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\n%s", RequestMenuTemplate, RequestMenuTextTemplate(data.Req)))
			keyboard, err := RequestKeyboard(stack.Data)
			if err != nil {
				data.Req = nil
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
				data.Req = nil
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
					log.Println("request back")
					data.Req = nil
					return ReturnOnParent(stack)
				}
			case "❌":
				err := stack.Bot.MentorService.UpdateRequest(context.Background(), &models.HelpRequest{
					ID:      data.Req.ID,
					GroupID: data.Group.ID,
					Status:  "rejected",
				})
				data.Req.Status = "rejected"
				if err != nil {
					data.Req = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nВы уже приняли заявку от этого студента!", ErrorMenuTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					} else if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nЗапрос не найден!", ErrorMenuTemplate)))
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

				id, err := stack.Bot.UsersService.GetTelegramID(context.Background(), data.Req.UserID)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Ваш запрос ментору %s, с целью %s был отклонён 😢", data.Req.Mentor.Name, data.Req.Goal)))
				if err != nil {
					log.Println(err)
				}
				return ReturnOnParent(stack)
			case "✔️":
				err := stack.Bot.MentorService.UpdateRequest(context.Background(), &models.HelpRequest{
					ID:      data.Req.ID,
					GroupID: data.Group.ID,
					Status:  "accepted",
				})
				data.Req.Status = "accepted"
				if err != nil {
					data.Req = nil
					data.LastMes = -1
					if err.(*httpError.HTTPError).StatusCode == http.StatusConflict {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nВы уже приняли заявку от этого студента!", ErrorMenuTemplate)))
						if err != nil {
							log.Println(err)
							return ReturnOnParent(stack)
						}
					} else if err.(*httpError.HTTPError).StatusCode == http.StatusNotFound {
						_, err := stack.Bot.Api.Send(tgbotapi.NewMessage(stack.ChatID, fmt.Sprintf("%s\n\nЗапрос не найден!", ErrorMenuTemplate)))
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
				id, err := stack.Bot.UsersService.GetTelegramID(context.Background(), data.Req.UserID)
				if err != nil {
					log.Println(err)
					return ReturnOnParent(stack)
				}
				_, err = stack.Bot.Api.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Ваш запрос ментору %s, с целью %s был принят 🤩", data.Req.Mentor.Name, data.Req.Goal)))
				if err != nil {
					log.Println(err)
				}
				return ReturnOnParent(stack)
			default:
				data.LastMes = -1
				return stack
			}
		}
	}
	return stack
}

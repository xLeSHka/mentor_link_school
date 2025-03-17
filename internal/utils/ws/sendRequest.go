package ws

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"log"
)

func SendRequest(personId, mentorId, requestId, groupId uuid.UUID, producer *broker.Producer, usersService service.UsersService, minioRepository repository.MinioRepository, studentService service.StudentService) {
	if producer != nil {
		user, err := usersService.GetByID(context.Background(), personId)
		if err != nil {
			log.Println(err)
			return
		}
		err = avatar.GetUserAvatar(user, minioRepository)
		if err != nil {
			log.Println(err)
			return
		}
		mentor, err := usersService.GetByID(context.Background(), mentorId)
		if err != nil {
			log.Println(err)
			return
		}
		err = avatar.GetUserAvatar(mentor, minioRepository)
		if err != nil {
			log.Println(err)
			return
		}
		request, err := studentService.GetRequestByID(context.Background(), requestId, groupId)
		if err != nil {
			log.Println(err)
			return
		}
		msg := &ws.Message{
			Type:       "request",
			UserID:     personId,
			TelegramID: user.TelegramID,
			Request: &ws.Request{
				ID:              requestId,
				StudentID:       personId,
				MentorID:        mentorId,
				MentorName:      mentor.Name,
				StudentName:     user.Name,
				MentorUrl:       mentor.AvatarURL,
				StudentUrl:      user.AvatarURL,
				StudentBio:      user.BIO,
				MentorBio:       mentor.BIO,
				StudentTelegram: user.Telegram,
				MentorTelegram:  mentor.Telegram,
				Goal:            request.Goal,
				Status:          request.Status,
			},
		}
		jsondData, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}
		err = producer.Send(jsondData)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

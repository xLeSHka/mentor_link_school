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

func SendUser(personId uuid.UUID, producer *broker.Producer, usersService service.UsersService, minioRepository repository.MinioRepository) {
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
		msg := &ws.Message{
			Type:   "user",
			UserID: personId,
			User: &ws.User{
				Name:      user.Name,
				AvatarUrl: user.AvatarURL,
				Telegram:  user.Telegram,
				BIO:       user.BIO,
			},
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}

		err = producer.Send(jsonData)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

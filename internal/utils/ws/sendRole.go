package ws

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/broker"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/avatar"
	"log"
)

func SendRole(personId, groupID uuid.UUID, role string, producer *broker.Producer, minioRepository repository.MinioRepository, groupService service.GroupService) {
	if producer != nil {
		group, err := groupService.GetGroupByID(context.Background(), groupID)
		if err != nil {
			log.Println(err)
			return
		}
		if group.AvatarURL != nil {
			err = avatar.GetGroupAvatar(group, minioRepository)
			if err != nil {
				log.Println(err)
				return
			}
		}
		mes := &ws.Message{
			Type:   "role",
			UserID: personId,
			Role: &ws.Role{
				Role:     role,
				GroupID:  groupID,
				GroupUrl: group.AvatarURL,
				Name:     group.Name,
			},
		}
		if role == "owner" {
			mes.Role.InviteCode = group.InviteCode
		}
		err = producer.Send(mes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

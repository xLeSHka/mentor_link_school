package groupsRoute

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

type reqGetMentorDto struct {
	GroupEmail string `json:"group_email" binding:"required,min=8,max=120,email"`
}
type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type reqCreateGroupDto struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type respCreateGroup struct {
	GroupID uuid.UUID `json:"group_id"`
}
type respGetGroupDto struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	AvatarUrl  *string `json:"avatar_url,omitempty"`
	InviteCode *string `json:"invite_code,omitempty"`
}

func mapGroup(group *models.Group) *respGetGroupDto {
	return &respGetGroupDto{
		Name:       group.Name,
		ID:         group.ID.String(),
		AvatarUrl:  group.AvatarURL,
		InviteCode: group.InviteCode,
	}
}

type resGetMember struct {
	Name      string   `json:"name"`
	AvatarUrl *string  `json:"avatar_url,omitempty"`
	BIO       *string  `json:"bio,omitempty"`
	Roles     []string `json:"roles"`
}

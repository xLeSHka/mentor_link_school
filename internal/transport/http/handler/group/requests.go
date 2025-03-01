package groupsRoute

import "gitlab.prodcontest.ru/team-14/lotti/internal/models"

type reqGetMentorDto struct {
	GroupEmail string `json:"group_email" binding:"required,min=8,max=120,email"`
}
type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type reqCreateGroupDto struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,min=8,max=120,email"`
}
type respGetGroupDto struct {
	Name      string  `json:"name"`
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
}

func mapGroup(group *models.Group) *respGetGroupDto {
	return &respGetGroupDto{
		Name: group.Name,
		ID:   group.ID.String(),
		//Email:     group.Email,
		AvatarUrl: group.AvatarURL,
	}
}

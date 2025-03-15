package groupsRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type ReqEditGroup struct {
	Name string `json:"name" binding:"required"`
}
type reqGetMentorDto struct {
	GroupEmail string `json:"group_email" binding:"required,min=8,max=120,email"`
}
type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type ReqCreateGroupDto struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}
type ReqUpdateRole struct {
	Role string `json:"role" binding:"required"`
	ID   string `json:"id" binding:"required,uuid"`
}
type GespGetMember struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Role      string    `uri:"role" binding:"required"`
}

func mapMember(role *models.Role) *GespGetMember {
	return &GespGetMember{
		UserID:    role.User.ID,
		AvatarUrl: role.User.AvatarURL,
		Name:      role.User.Name,
		Role:      role.Role,
	}
}

type respUploadAvatarDto struct {
	Url string `json:"url"`
}

type respJoinGroup struct {
	Status string `json:"status"`
}

type respStat struct {
	StudentsCount        int64   `json:"students_count"`
	MentorsCount         int64   `json:"mentors_count"`
	HelpRequestCount     int64   `json:"help_request_count"`
	AcceptedRequestCount int64   `json:"accepted_request_count"`
	RejectedRequestCount int64   `json:"rejected_request_count"`
	Conversion           float64 `json:"conversion"`
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

//type resGetMember struct {
//	Name      string  `json:"name"`
//	AvatarUrl *string `json:"avatar_url,omitempty"`
//	BIO       *string `json:"bio,omitempty"`
//	Role      string  `json:"role"`
//}

type reqUpdateRoleDto struct {
	MemberID uuid.UUID `json:"member_id" binding:"required"`
	Roles    string    `json:"roles" binding:"required"`
}

type respUpdateCode struct {
	Code string `json:"code"`
}

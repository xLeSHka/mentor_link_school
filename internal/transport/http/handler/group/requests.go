package groupsRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type ReqEditGroup struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
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
	Role string `json:"role" binding:"required" validate:"c-role"`
}
type RespGetMember struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Roles     []string  `json:"roles" binding:"required"`
}
type OffsetRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func mapMember(role *models.User) *RespGetMember {
	roles := []string{}
	for _, role := range role.Roles {
		roles = append(roles, role.Role)
	}
	return &RespGetMember{
		UserID:    role.ID,
		AvatarUrl: role.AvatarURL,
		Name:      role.Name,
		Roles:     roles,
	}
}

type RespUploadAvatarDto struct {
	Url string `json:"url"`
}

type RespStat struct {
	StudentsCount        int64   `json:"students_count"`
	MentorsCount         int64   `json:"mentors_count"`
	HelpRequestCount     int64   `json:"help_request_count"`
	AcceptedRequestCount int64   `json:"accepted_request_count"`
	RejectedRequestCount int64   `json:"rejected_request_count"`
	Conversion           float64 `json:"conversion"`
}
type RespCreateGroup struct {
	GroupID    uuid.UUID `json:"group_id"`
	InviteCode string    `json:"invite_code"`
}

//type respGetGroupDto struct {
//	Name       string  `json:"name"`
//	ID         string  `json:"id"`
//	AvatarUrl  *string `json:"avatar_url,omitempty"`
//	InviteCode *string `json:"invite_code,omitempty"`
//}
//
//func MapGroup(group *models.Group) *respGetGroupDto {
//	return &respGetGroupDto{
//		Name:       group.Name,
//		ID:         group.ID.String(),
//		AvatarUrl:  group.AvatarURL,
//		InviteCode: group.InviteCode,
//	}
//}

//type resGetMember struct {
//	Name      string  `json:"name"`
//	AvatarUrl *string `json:"avatar_url,omitempty"`
//	BIO       *string `json:"bio,omitempty"`
//	Role      string  `json:"role"`
//}

type RespUpdateCode struct {
	Code string `json:"code"`
}
type RespGetRoles struct {
	Roles []string `json:"roles"`
}

func MapRoles(roles []*models.Role) *RespGetRoles {
	resp := []string{}
	for _, role := range roles {
		resp = append(resp, role.Role)
	}
	return &RespGetRoles{Roles: resp}
}

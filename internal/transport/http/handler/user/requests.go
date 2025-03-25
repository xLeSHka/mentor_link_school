package usersRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type OffsetRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type ReqRegisterDto struct {
	Name     *string `form:"name" binding:"required,gte=1,lte=120"`
	Telegram *string `json:"telegram" binding:"required,gte=1,lte=120"`
	Password *string `json:"password" binding:"required,gte=8,lte=60" validate:"c-password"`
}
type RespRegisterDto struct {
	Token string `json:"token"`
}
type ReqLoginDto struct {
	Telegram *string `json:"telegram" binding:"required,gte=1,lte=120"`
	Password *string `json:"password" binding:"required,gte=8,lte=60" validate:"c-password"`
}
type RespLoginDto struct {
	Token string `json:"token"`
}
type ReqOtherProfileDto struct {
	ProfileID *string `uri:"profileID" binding:"required,uuid"`
}

type ResGetGroup struct {
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	Roles      []string  `json:"roles"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	InviteCode *string   `json:"invite_code,omitempty"`
}
type ResGetProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
}
type RespUploadAvatarDto struct {
	Url string `json:"url"`
}
type RespOtherProfile struct {
	ID        uuid.UUID `json:"id"`
	Telegram  string    `json:"telegram"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	BIO       *string   `json:"bio,omitempty"`
}
type ReqEditUser struct {
	Telegram *string `form:"telegram" binding:"omitempty,gte=1,lte=120"`
	Name     *string `json:"name" binding:"omitempty,gte=1,lte=120"`
	BIO      *string `json:"bio" binding:"omitempty,lte=500"`
}
type RespJoinGroup struct {
	Status string `json:"status"`
}

func MapOtherProfile(user *models.User) *RespOtherProfile {
	return &RespOtherProfile{
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Telegram:  user.Telegram,
		BIO:       user.BIO,
	}
}
func MapProfile(user *models.User) *ResGetProfile {
	return &ResGetProfile{
		ID:        user.ID,
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		Telegram:  user.Telegram,
		BIO:       user.BIO,
	}
}

func MapGroup(roles *models.GroupWithRoles) *ResGetGroup {
	resp := &ResGetGroup{
		GroupID:   roles.GroupID,
		Name:      roles.Group.Name,
		Roles:     roles.MyRoles,
		AvatarURL: roles.Group.AvatarURL,
	}
	for _, role := range roles.MyRoles {
		if role == "owner" {
			resp.InviteCode = roles.Group.InviteCode
		}
	}
	return resp
}

package usersRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

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
	ProfileID *string `uri:"id" binding:"required,uuid"`
}
type reqGetRole struct {
	Role string `from:"role" binding:"required"`
}
type Role struct {
	GroupID uuid.UUID `json:"group_id"`
	Role    string    `json:"role"`
}
type ResGetProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
	Roles     []*Role   `json:"roles"`
}

type RespGetMyMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Telegram  string    `json:"telegram"`
	BIO       *string   `json:"bio,omitempty"`
}
type RespGetMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
}
type RespGetHelp struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	MentorName string    `json:"mentor_name"`
	AvatarUrl  *string   `json:"avatar_url,omitempty"`
	Goal       string    `json:"goal"`
	Status     string    `json:"status"`
	Telegram   string    `json:"mentor_telegram"`
	BIO        *string   `json:"mentor_bio,omitempty"`
}
type Pair struct {
	MentorID uuid.UUID `json:"mentor_id"`
	GroupId  uuid.UUID `json:"group_id"`
}
type ReqCreateHelp struct {
	Requests []Pair `json:"requests" binding:"required"`
	Goal     string `json:"goal" binding:"required"`
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
	roles := make([]*Role, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, MapRole(role))
	}
	return &ResGetProfile{
		ID:        user.ID,
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		Telegram:  user.Telegram,
		BIO:       user.BIO,
		Roles:     roles,
	}
}

func MapMyMentor(mentor *models.Pair) *RespGetMyMentor {
	return &RespGetMyMentor{
		MentorID:  mentor.Mentor.ID,
		AvatarUrl: mentor.Mentor.AvatarURL,
		Name:      mentor.Mentor.Name,
		Telegram:  mentor.Mentor.Telegram,
		BIO:       mentor.Mentor.BIO,
	}
}
func MapHelp(help *models.HelpRequest) *RespGetHelp {
	return &RespGetHelp{
		ID:         help.ID,
		MentorID:   help.MentorID,
		Status:     help.Status,
		Goal:       help.Goal,
		MentorName: help.Mentor.Name,
		AvatarUrl:  help.Student.AvatarURL,
		Telegram:   help.Mentor.Telegram,
		BIO:        help.Mentor.BIO,
	}
}
func MapMentor(mentor *models.Role) *RespGetMentor {
	return &RespGetMentor{
		MentorID:  mentor.User.ID,
		AvatarUrl: mentor.User.AvatarURL,
		Name:      mentor.User.Name,
		BIO:       mentor.User.BIO,
		Telegram:  mentor.User.Telegram,
	}
}

func MapRole(role *models.Role) *Role {
	resp := &Role{
		GroupID: role.GroupID,
		Role:    role.Role,
	}
	return resp
}

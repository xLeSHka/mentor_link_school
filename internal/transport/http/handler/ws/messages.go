package ws

import (
	"github.com/google/uuid"
)

type Role struct {
	Role       string    `json:"role"`
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	GroupUrl   *string   `json:"group_url,omitempty"`
	InviteCode *string   `json:"invite_code,omitempty"`
}

//type RespGetGroupDto struct {
//	Name       string  `json:"name"`
//	ID         string  `json:"id"`
//	AvatarUrl  *string `json:"avatar_url,omitempty"`
//	InviteCode *string `json:"invite_code,omitempty"`
//	Role       string  `json:"role"`
//}
//
//func MapGroup(group *models.Group, role string) *RespGetGroupDto {
//	resp := &RespGetGroupDto{
//		Name:      group.Name,
//		ID:        group.ID.String(),
//		AvatarUrl: group.AvatarURL,
//		Role:      role,
//	}
//	if role == "owner" {
//		resp.InviteCode = group.InviteCode
//	}
//	return resp
//}

type User struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	Telegram  string  `json:"telegram"`
	BIO       *string `json:"bio,omitempty"`
}
type Request struct {
	ID              uuid.UUID `json:"id"`
	StudentID       uuid.UUID `json:"student_id"`
	MentorID        uuid.UUID `json:"mentor_id"`
	MentorName      string    `json:"mentor_name"`
	StudentName     string    `json:"student_name"`
	MentorUrl       *string   `json:"mentor_url,omitempty"`
	StudentUrl      *string   `json:"student_url,omitempty"`
	StudentBio      *string   `json:"student_bio,omitempty"`
	MentorBio       *string   `json:"mentor_bio,omitempty"`
	StudentTelegram string    `json:"student_telegram"`
	MentorTelegram  string    `json:"mentor_telegram"`
	Goal            string    `json:"goal"`
	Status          string    `json:"status"`
}
type Message struct {
	Type    string    `json:"type"`
	UserID  uuid.UUID `json:"user_id"`
	Role    *Role     `json:"role,omitempty"`
	User    *User     `json:"user,omitempty"`
	Request *Request  `json:"request,omitempty"`
}

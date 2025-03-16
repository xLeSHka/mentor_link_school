package studentsRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
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
type ReqCreateHelp struct {
	Goal string `json:"goal" binding:"required"`
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

package usersRoute

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

type reqLoginDto struct {
	Name string `json:"name" binding:"required"`
}
type respLoginDto struct {
	Token string `json:"token"`
}

type resGetProfile struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}

type respGetMyMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}
type respGetMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupID   uuid.UUID `json:"group_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
}
type respGetHelp struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	MentorName string    `json:"mentor_name"`
	AvatarUrl  *string   `json:"avatar_url,omitempty"`
	Goal       string    `json:"goal"`
	Status     string    `json:"status"`
}
type Pair struct {
	MentorID uuid.UUID `json:"mentor_id"`
	GroupId  uuid.UUID `json:"group_id"`
}
type reqCreateHelp struct {
	Requests []Pair `json:"requests" binding:"required"`
	Goal     string `json:"goal" binding:"required"`
}
type respUploadAvatarDto struct {
	Url string `json:"url"`
}

func mapMyMentor(mentor *models.Pair) *respGetMyMentor {
	return &respGetMyMentor{
		MentorID:  mentor.Mentor.ID,
		AvatarUrl: mentor.Mentor.AvatarURL,
		Name:      mentor.Mentor.Name,
	}
}
func mapHelp(help *models.HelpRequest) *respGetHelp {
	return &respGetHelp{
		ID:         help.ID,
		MentorID:   help.MentorID,
		Status:     help.Status,
		Goal:       help.Goal,
		MentorName: help.Mentor.Name,
		AvatarUrl:  help.Student.AvatarURL,
	}
}
func mapMentor(mentor *models.Role) *respGetMentor {
	return &respGetMentor{
		MentorID:  mentor.User.ID,
		GroupID:   mentor.GroupID,
		AvatarUrl: mentor.User.AvatarURL,
		Name:      mentor.User.Name,
		BIO:       mentor.User.BIO,
	}
}

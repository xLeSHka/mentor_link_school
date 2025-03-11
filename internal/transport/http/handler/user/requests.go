package usersRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type reqLoginDto struct {
	Name string `json:"name" binding:"required"`
}
type respLoginDto struct {
	Token string `json:"token"`
}

type reqGetRole struct {
	Role string `from:"role" binding:"required"`
}
type resGetInitData struct {
	Name      string             `json:"name"`
	AvatarUrl *string            `json:"avatar_url,omitempty"`
	BIO       *string            `json:"bio,omitempty"`
	Telegram  *string            `json:"telegram"`
	Groups    []*RespGetGroupDto `json:"groups"`
}

type respGetMyMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupIDs  []string  `json:"group_ids" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Telegram  string    `json:"telegram"`
	BIO       *string   `json:"bio,omitempty"`
}
type respGetMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupIDs  []string  `json:"group_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
}
type respGetHelp struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	GroupIDs   []string  `json:"group_ids"`
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
type reqCreateHelp struct {
	Requests []Pair `json:"requests" binding:"required"`
	Goal     string `json:"goal" binding:"required"`
}
type respUploadAvatarDto struct {
	Url string `json:"url"`
}
type respOtherProfile struct {
	Telegram string  `json:"telegram"`
	BIO      *string `json:"bio,omitempty"`
}
type reqEditUser struct {
	Name     string `json:"name" binding:"required"`
	Telegram string `json:"telegram,required"`
	BIO      string `json:"bio,required"`
}

func mapOtherProfile(user *models.User) *respOtherProfile {
	return &respOtherProfile{
		Telegram: user.Telegram,
		BIO:      user.BIO,
	}
}
func mapMyMentor(mentor *models.PairWithGIDs) *respGetMyMentor {
	return &respGetMyMentor{
		MentorID:  mentor.Mentor.ID,
		AvatarUrl: mentor.Mentor.AvatarURL,
		GroupIDs:  mentor.GroupIDs,
		Name:      mentor.Mentor.Name,
		Telegram:  mentor.Mentor.Telegram,
		BIO:       mentor.Mentor.BIO,
	}
}
func mapHelp(help *models.HelpRequestWithGIDs) *respGetHelp {
	return &respGetHelp{
		ID:         help.ID,
		MentorID:   help.MentorID,
		GroupIDs:   help.GroupIDs,
		Status:     help.Status,
		Goal:       help.Goal,
		MentorName: help.Mentor.Name,
		AvatarUrl:  help.Student.AvatarURL,
		Telegram:   help.Mentor.Telegram,
		BIO:        help.Mentor.BIO,
	}
}
func mapMentor(mentor *models.RoleWithGIDs) *respGetMentor {
	return &respGetMentor{
		MentorID:  mentor.User.ID,
		GroupIDs:  mentor.GroupIDs,
		AvatarUrl: mentor.User.AvatarURL,
		Name:      mentor.User.Name,
		BIO:       mentor.User.BIO,
		Telegram:  mentor.User.Telegram,
	}
}

type RespGetGroupDto struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	AvatarUrl  *string `json:"avatar_url,omitempty"`
	InviteCode *string `json:"invite_code,omitempty"`
	Role       string  `json:"role"`
}

func mapGroup(group *models.Group, role string) *RespGetGroupDto {
	resp := &RespGetGroupDto{
		Name:      group.Name,
		ID:        group.ID.String(),
		AvatarUrl: group.AvatarURL,
		Role:      role,
	}
	if role == "owner" {
		resp.InviteCode = group.InviteCode
	}
	return resp
}

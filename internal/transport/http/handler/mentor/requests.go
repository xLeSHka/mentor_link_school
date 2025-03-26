package mentorsRoute

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

type OffsetRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type ReqUpdateRequest struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status *bool     `json:"status" binding:"required"`
}
type RespGetMyStudent struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
}
type RespGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
	Telegram  string    `json:"student_telegram"`
	BIO       *string   `json:"student_bio,omitempty"`
}

func MapMyStudent(user *models.Pair) *RespGetMyStudent {
	return &RespGetMyStudent{
		StudentID: user.UserID,
		AvatarUrl: user.Student.AvatarURL,
		Name:      user.Student.Name,
		BIO:       user.Student.BIO,
		Telegram:  user.Student.Telegram,
	}
}
func MapRequest(req *models.HelpRequest) (res *RespGetRequest) {
	res = &RespGetRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		Name:      req.Student.Name,
		AvatarUrl: req.Student.AvatarURL,
		Goal:      req.Goal,
		Status:    req.Status,
		Telegram:  req.Student.Telegram,
		BIO:       req.Student.BIO,
	}

	return
}

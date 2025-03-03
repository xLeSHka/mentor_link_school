package mentorsRoute

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

type reqUpdateRequest struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status *bool     `json:"status" binding:"required"`
}
type respGetMyStudent struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	GroupIDs  []string  `json:"group_ids" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
	Telegram  string    `json:"telegram"`
}
type respGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	GroupIDs  []string  `json:"group_ids"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
}

func mapMyStudent(user *models.PairWithGIDs) *respGetMyStudent {
	return &respGetMyStudent{
		StudentID: user.UserID,
		GroupIDs:  user.GroupIDs,
		AvatarUrl: user.Student.AvatarURL,
		Name:      user.Student.Name,
		BIO:       user.Student.BIO,
		Telegram:  user.Student.Telegram,
	}
}
func mapRequest(req *models.HelpRequestWithGIDs) (res *respGetRequest) {
	res = &respGetRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		GroupIDs:  req.GroupIDs,
		Name:      req.Student.Name,
		AvatarUrl: req.Student.AvatarURL,
		Goal:      req.Goal,
		Status:    req.Status,
	}

	return
}

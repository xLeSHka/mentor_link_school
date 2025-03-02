package mentorsRoute

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

type reqUpdateRequest struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status bool      `json:"status" binding:"required"`
}
type respGetMyStudent struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}
type respGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
}

func mapMyStudent(user *models.User) *respGetMyStudent {
	return &respGetMyStudent{
		StudentID: user.ID,
		AvatarUrl: user.AvatarURL,
		Name:      user.Name,
	}
}
func mapRequest(req *models.HelpRequest) (res *respGetRequest) {
	res = &respGetRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		Name:      req.Student.Name,
		AvatarUrl: req.Student.AvatarURL,
		Goal:      req.Goal,
		Status:    req.Status,
	}

	return
}

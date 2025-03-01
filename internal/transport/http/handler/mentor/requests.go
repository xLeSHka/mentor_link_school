package mentorsRoute

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

type reqUpdateRequest struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status bool      `json:"status" binding:"required"`
}
type resGetProfile struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Name      string    `json:"name"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
}
type respGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Goal      string    `json:"goal"`
}

func mapRequest(req *models.HelpRequest) *respGetRequest {
	return &respGetRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		AvatarURL: req.Student.AvatarURL,
		Name:      req.Student.Name,
		Goal:      req.Goal,
	}
}

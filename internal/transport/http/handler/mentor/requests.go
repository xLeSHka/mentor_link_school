package mentorsRoute

import "github.com/google/uuid"

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
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Goal      string    `json:"goal"`
}

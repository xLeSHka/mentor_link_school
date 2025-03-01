package repositoryUser

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *UsersRepository) GetMyRequests(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error) {
	var resp []*models.HelpRequest
	res := r.DB.Model(&models.HelpRequest{}).Where("user_id = ? AND status = 'pending'", userID).
		Preload("Mentor", "users.id = help_requests.mentor_id").
		Joins("JOIN users ON users.id = help_requests.mentor_id").
		Find(&resp)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return resp, nil
		}
		return resp, res.Error
	}
	return resp, nil
}

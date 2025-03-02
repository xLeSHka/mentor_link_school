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
		WithContext(ctx).
		Preload("User").
		Preload("Student").
		Find(&resp)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return resp, nil
		}
		return resp, res.Error
	}
	return resp, nil
}

package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetMyMentors(ctx context.Context, id uuid.UUID) ([]*models.Pair, error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).Where("user_id = ?", id).
		WithContext(ctx).
		Preload("User").
		Joins("JOIN users ON users.id = pairs.mentor_id").
		Find(&resp).Error
	return resp, err
}

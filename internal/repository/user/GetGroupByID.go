package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetGroupByID(ctx context.Context, ID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.DB.WithContext(ctx).First(&group, &models.Group{ID: ID}).Error
	return &group, err
}

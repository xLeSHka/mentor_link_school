package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

// func (r *GroupRepository) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error) {
// 	var groups []*models.Group
// 	err := r.DB.Model(&models.Group{}).WithContext(ctx).Where("user_id = ?", userID).Find(&groups).Error
// 	return groups, err
// }

func (r *UsersRepository) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	var groups []*models.Role
	err := r.DB.WithContext(ctx).Table("roles").
		Preload("Group").
		Where("user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

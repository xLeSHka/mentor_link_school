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

func (r *UsersRepository) GetGroups(ctx context.Context, userID uuid.UUID, role string) ([]*models.Group, error) {
	var groupIDs []string

	err := r.DB.Table("roles").Select("group_id").Where("user_id = ? AND role = ?", userID, role).Find(&groupIDs).Error
	if err != nil {
		return nil, err
	}
	var groups []*models.Group
	err = r.DB.Model(&models.Group{}).Where("id in (?)", groupIDs).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, err
}

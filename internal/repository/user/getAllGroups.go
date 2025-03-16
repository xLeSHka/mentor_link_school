package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

// func (r *GroupRepository) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error) {
// 	var groups []*models.Group
// 	err := r.DB.Model(&models.Group{}).WithContext(ctx).Where("user_id = ?", userID).Find(&groups).Error
// 	return groups, err
// }

func (r *UsersRepository) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Roles, error) {
	var groups []*models.Roles
	err := r.DB.WithContext(ctx).Table("roles").
		Select("user_id,group_id,array_agg(role) AS rls").Group("group_id").Group("user_id").
		Preload("Group").
		Where("user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

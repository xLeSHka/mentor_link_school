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

func (r *UsersRepository) GetGroups(ctx context.Context, userID uuid.UUID, page, size int) ([]*models.GroupWithRoles, int64, error) {
	var groups []*models.GroupWithRoles
	err := r.DB.WithContext(ctx).Table("roles").
		Select("user_id,group_id,array_agg(role) AS my_roles").Group("group_id").Group("user_id").
		Preload("Group").
		Where("user_id = ?", userID).
		Offset(page * size).
		Limit(size).
		Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = r.DB.WithContext(ctx).Model(&models.Role{}).Group("group_id").Where("user_id = ?", userID).Count(&total).Error
	return groups, total, err
}

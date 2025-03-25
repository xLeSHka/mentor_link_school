package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *GroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID, page, size int) ([]*models.User, int64, error) {
	var members []*models.User
	err := r.DB.Model(&models.User{}).Where("EXISTS (SELECT 1 FROM roles WHERE group_id = ? AND user_id = users.id)", groupID).
		Preload("Roles", "roles.group_id = ? ", groupID).
		Offset(page * size).
		Limit(size).
		Find(&members).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = r.DB.Model(&models.Role{}).Group("user_id").Where("group_id = ?", groupID).Count(&total).Error
	return members, total, nil
}

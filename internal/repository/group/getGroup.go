package group

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) GetGroup(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.DB.WithContext(ctx).
		Table("groups").
		Joins("JOIN roles ON roles.group_id = groups.id").
		Where("roles.user_id = ? AND roles.role = 'owner' AND groups.id = ?", userID, groupID).
		Select("groups.*").
		First(&group).Error

	if err != nil {
		return nil, err
	}

	return &group, nil
}

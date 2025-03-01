package group

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) GetGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	var resp models.Group
	err := r.DB.Model(&models.Group{}).WithContext(ctx).Where("id = ? AND user_id = ?", group.ID, group.UserID).First(&resp).Error
	return &resp, err
}

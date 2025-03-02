package group

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) GetGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*models.Group, error) {
	var resp models.Group
	err := r.DB.Model(&models.Role{}).WithContext(ctx).Where("user_id = ? AND role = 'owner' AND group_id = ?", userID, groupID).First(&resp).Error
	return &resp, err
}

package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *GroupRepository) GetRoles(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error) {
	var members []*models.Role
	err := r.DB.Model(&models.Role{}).Where("user_id = ? AND group_id = ?", userID, groupID).
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

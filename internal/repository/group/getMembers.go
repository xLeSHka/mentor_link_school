package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *GroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.Role, error) {
	var members []*models.Role
	err := r.DB.Model(&models.Role{}).Where("group_id = ? AND role != 'owner'", groupID).
		Preload("User").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

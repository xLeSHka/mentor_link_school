package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *GroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.User, error) {
	var members []*models.User
	err := r.DB.Model(&models.User{}).Where("EXISTS (SELECT 1 FROM roles WHERE group_id = ? AND user_id = users.id)", groupID).
		Preload("Rls").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

package group

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) GetAllCampaigns(userID uuid.UUID) ([]*models.Group, error) {
	var groups []*models.Group
	err := r.DB.Model(&models.Group{}).Where("user_id = ?", userID).Find(&groups).Error
	return groups, err
}

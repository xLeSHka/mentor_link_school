package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (r *GroupRepository) CheckGroupExists(ctx context.Context, user_id, group_id uuid.UUID) (bool, error) {
	var group models.Role
	err := r.DB.Model(&models.Role{}).Where("user_id = ? AND group_id = ? AND role = 'owner'", user_id, group_id).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

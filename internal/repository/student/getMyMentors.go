package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetMyMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).WithContext(ctx).
		Where("user_id = ? AND group_id = ?", userID, groupID).Preload("Mentor").Find(&resp).Error
	return resp, err
}

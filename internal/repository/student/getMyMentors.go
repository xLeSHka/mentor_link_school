package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetMyMentors(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.Pair, int64, error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).WithContext(ctx).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Preload("Mentor").
		Offset(page * size).
		Limit(size).
		Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = r.DB.Model(&models.Pair{}).Where("user_id = ? AND group_id = ?", userID, groupID).Count(&count).Error
	return resp, count, err
}

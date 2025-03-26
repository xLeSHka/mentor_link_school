package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetMyRequests(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.HelpRequest, int64, error) {
	var resp []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).
		Where("user_id = ? AND group_id = ? ", userID, groupID).Preload("Mentor").
		Preload("Student").Order("status").
		Offset(page * size).
		Limit(size).
		Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = r.DB.Model(&models.HelpRequest{}).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Count(&count).Error
	return resp, count, err
}

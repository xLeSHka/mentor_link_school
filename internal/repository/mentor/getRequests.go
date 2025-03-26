package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.HelpRequest, int64, error) {

	var resp []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).
		Where("mentor_id = ? AND status = 'pending' AND group_id = ?", userID, groupID).Preload("Student").
		Preload("Mentor").
		Offset(page * size).Limit(size).
		Find(&resp).Error
	if err != nil {
		return resp, 0, err
	}
	var count int64
	err = r.DB.Model(&models.HelpRequest{}).Where("mentor_id = ? AND status = 'pending' AND group_id = ?", userID, groupID).Count(&count).Error
	return resp, count, err
}

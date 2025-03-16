package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error) {

	var resp []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).
		Where("mentor_id = ? AND status = 'pending' AND group_id = ?", userID, groupID).Preload("Student").
		Preload("Mentor").
		Find(&resp).Error
	return resp, err
}

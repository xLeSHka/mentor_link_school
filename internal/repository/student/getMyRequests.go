package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetMyRequests(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error) {
	var resp []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).
		Where("user_id = ? AND group_id = ? AND (status = 'pending' OR status = 'rejected')", userID, groupID).Preload("Mentor").
		Where("NOT EXISTS (SELECT 1 FROM pairs WHERE user_id = ? and mentor_id = help_requests.mentor_id AND group_id = ?)", userID, groupID).
		Preload("Student").Order("status").Find(&resp).Error
	return resp, err
}

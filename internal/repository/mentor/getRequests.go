package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error) {
	var request []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).Where("mentor_id = ? AND status ='pending'", userID).
		Preload("Student", "users.id = help_requests.student_id").
		Joins("JOIN users ON users.id = help_requests.student_id").
		Find(&request).Error
	return request, err
}

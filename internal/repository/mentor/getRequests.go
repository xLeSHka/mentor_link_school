package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error) {
	var request []*models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).Where("mentor_id = ? AND status = 'pending'", userID).
		Preload("User").
		Preload("Student").
		Find(&request).Error
	return request, err
}

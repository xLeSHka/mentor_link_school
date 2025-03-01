package repositoryMentor

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).
		Where("id = ?", request.ID).
		Updates(map[string]any{"status": request.Status}).Error
	return err
}

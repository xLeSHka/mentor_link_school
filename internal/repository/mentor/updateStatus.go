package repositoryMentor

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	pair := &models.Pair{
		UserID:   request.UserID,
		MentorID: request.MentorID,
		Goal:     request.Goal,
	}

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pair).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE help_requests SET status = 'accepted' WHERE id = ?", request.ID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

package repositoryMentor

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	var reqData models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).
		WithContext(ctx).Where("id = ?", request.ID).
		First(&reqData).Error
	if err != nil {
		return err
	}

	pair := &models.Pair{
		GroupID:  reqData.GroupID,
		UserID:   reqData.UserID,
		MentorID: reqData.MentorID,
		Goal:     reqData.Goal,
	}

	err = r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

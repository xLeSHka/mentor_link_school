package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error) {
	var resp []*models.HelpRequestWithGIDs
	err := r.DB.Table("help_requests").
		Select("id,user_id,mentor_id,array_agg(group_id) as group_ids,goal,status").
		Where("mentor_id = ? AND status = 'pending'", userID).Group("id,mentor_id").Group("user_id").Group("goal").Preload("Student").
		Where("NOT EXISTS (SELECT 1 FROM pairs WHERE mentor_id = ? and user_id = help_requests.user_id)", userID).
		Preload("Mentor").
		Find(&resp).Error
	return resp, err
}

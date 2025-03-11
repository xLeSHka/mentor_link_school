package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) GetMyRequests(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error) {
	var resp []*models.HelpRequestWithGIDs
	err := r.DB.Table("help_requests").
		Select("id,user_id,mentor_id,array_agg(group_id) as group_ids,goal,status").
		Where("user_id = ? AND status = 'pending'", userID).Group("id,mentor_id").Group("user_id").Group("goal").Preload("Mentor").
		Where("NOT EXISTS (SELECT 1 FROM pairs WHERE user_id = ? and mentor_id = help_requests.mentor_id)", userID).
		Preload("Student").
		Find(&resp).Error
	return resp, err
}

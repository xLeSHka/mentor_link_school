package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error) {

	var resp []*models.HelpRequestWithGIDs
	err := r.DB.Table("help_requests").
		Select("id,user_id,mentor_id,array_agg(group_id) as group_ids,goal,status").Group("id,user_id,mentor_id").
		Where("mentor_id = ? AND status = 'pending'", userID).Preload("Student").
		Preload("Mentor").
		Find(&resp).Error
	return resp, err
}

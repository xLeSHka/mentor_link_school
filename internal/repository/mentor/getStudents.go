package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.PairWithGIDs, error) {
	var resp []*models.PairWithGIDs
	err := r.DB.Table("pairs").Preload("Student").Select("user_id,mentor_id,array_agg(group_id) as group_ids").Where("mentor_id = ?", userID).Group("mentor_id").Group("user_id").Find(&resp).Error
	return resp, err
}

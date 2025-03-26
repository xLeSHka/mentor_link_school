package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetStudents(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.Pair, int64, error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).Preload("Student").Where("mentor_id = ? AND group_id = ?", userID, groupID).
		Offset(page * size).Limit(size).
		Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = r.DB.Model(&models.Pair{}).Where("mentor_id = ? AND group_id = ?", userID, groupID).Count(&count).Error
	return resp, count, err
}

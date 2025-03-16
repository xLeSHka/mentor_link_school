package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) GetStudents(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).Preload("Student").Where("mentor_id = ? AND group_id = ?", userID, groupID).Find(&resp).Error
	return resp, err
}

package repositoryMentor

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (r *MentorRepository) CheckRequest(ctx context.Context, id, mentorID uuid.UUID) (bool, error) {
	var res models.HelpRequest
	err := r.DB.Model(&models.HelpRequest{}).WithContext(ctx).Where("id = ? AND mentor_id = ? AND status = 'pending'", id, mentorID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

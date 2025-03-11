package repositoryMentor

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (r *MentorRepository) CheckIsMentor(ctx context.Context, id, groupID uuid.UUID) (bool, error) {
	var res models.Role
	err := r.DB.Model(&models.Role{}).WithContext(ctx).Where("user_id = ? AND group_id = ? AND role = 'mentor' OR role = 'student-mentor'", id, groupID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

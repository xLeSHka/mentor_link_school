package repositoryMentor

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	result := r.DB.WithContext(ctx).Model(&models.HelpRequest{}).Where("id = ? AND group_id = ?", request.ID, request.GroupID).Update("status", request.Status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

package repositoryMentor

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	err := r.DB.WithContext(ctx).Model(&models.HelpRequest{}).Where("id = ? AND group_id = ?", request.ID, request.GroupID).Update("status", request.Status).Error
	return err
}

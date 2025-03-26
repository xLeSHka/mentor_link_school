package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetRequestByID(ctx context.Context, reqID, groupID uuid.UUID) (*models.HelpRequest, error) {
	var res models.HelpRequest
	err := r.DB.Model(models.HelpRequest{}).
		WithContext(ctx).
		Where("id = ? AND group_id = ?", reqID, groupID).Preload("Student").Preload("Mentor").
		First(&res).Error

	return &res, err
}

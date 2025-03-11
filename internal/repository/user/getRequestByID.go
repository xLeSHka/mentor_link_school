package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) GetRequestByID(ctx context.Context, reqID uuid.UUID) (models.HelpRequest, error) {
	var res models.HelpRequest
	err := r.DB.Model(models.HelpRequest{}).
		WithContext(ctx).
		Where("id = ?", reqID).
		First(&res).Error

	return res, err
}

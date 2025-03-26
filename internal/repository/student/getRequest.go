package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) GetRequest(ctx context.Context, studentID, mentorID, groupID uuid.UUID) (*models.HelpRequest, error) {
	var res models.HelpRequest
	err := r.DB.Model(models.HelpRequest{}).
		WithContext(ctx).
		Where("user_id = ? AND mentor_id = ? AND group_id = ? and (status = 'pending' OR status = 'accepted')", studentID, mentorID, groupID).
		First(&res).Error
	return &res, err
}

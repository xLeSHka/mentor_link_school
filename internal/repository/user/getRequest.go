package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetRequest(ctx context.Context, UserID, MentorID, GroupID uuid.UUID) (models.HelpRequest, error) {
	var res models.HelpRequest
	err := r.DB.Model(models.HelpRequest{}).
		WithContext(ctx).
		Where("user_id = ? AND mentor_id = ? AND group_id = ? and status = 'pending'", UserID, MentorID, GroupID).
		First(&res).Error

	return res, err
}

package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetMyMentors(ctx context.Context, id uuid.UUID) ([]*models.PairWithGIDs, error) {
	var resp []*models.PairWithGIDs
	err := r.DB.Table("pairs").Select("user_id,mentor_id,array_agg(group_id) as group_id").Where("user_id = ?", id).Group("mentor_id").Group("user_id").Preload("Mentor").Find(&resp).Error
	return resp, err
}

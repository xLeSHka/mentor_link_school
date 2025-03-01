package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetMyMentors(ctx context.Context, id uuid.UUID) (person *models.Pair, err error) {
	var resp []*models.Pair
	err := r.DB.Model(&models.Pair{}).Where("user_id = ?", id).
		Preload("Mentor", "mentors.id = pairs.mentor_id").
		Joins("JOIN mentors ON mentors.id = pairs.mentor.id").
		Find(&resp).Error
	return resp, err
}

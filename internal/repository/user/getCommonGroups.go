package repositoryUser

import (
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetCommonGroups(userID, mentorID uuid.UUID) ([]uuid.UUID, error) {
	var pairs []models.Pair
	err := r.DB.Model(&models.Pair{}).
		Where("user_id = ? AND mentor_id = ?", userID, mentorID).
		Find(&pairs).Error

	var ids []uuid.UUID
	for _, pair := range pairs {
		ids = append(ids, pair.GroupID)
	}

	return ids, err
}

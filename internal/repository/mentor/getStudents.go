package repositoryMentor

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *MentorRepository) GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error) {
	var users []*models.Pair
	err := r.DB.Model(&models.Pair{}).Where("mentor_id = ?", userID).
		Preload("Student").
		Joins("JOIN users ON users.id = pairs.user_id").
		Find(&users).Error
	return users, err
}

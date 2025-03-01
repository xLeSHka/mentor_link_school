package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) GetAvaliableMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error) {
	var resp []*models.Pair
	err := r.DB.
		Where("group_id in (select group_id from roles where user_id = ? AND role = 'student')", userID).
		Preload("Mentor").
		Table("users").Find(&resp).Error

	return resp, err

}

func (r *UsersRepository) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	var groups []string
	tx := r.DB.Begin()
	err := tx.Table("roles").Select("group_id").Where("user_id = ? AND role = 'student'", userID).Find(&groups).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var users []*models.User
	err = tx.Table("users").Select("*").
		Preload("Role", "role = 'mentor' AND group_id IN (?)", groups).
		Joins("CROSS JOIN roles").
		Where("group_id IN (?) AND role = 'mentor' AND id = user_id", groups).
		Find(&users).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return users, nil
}

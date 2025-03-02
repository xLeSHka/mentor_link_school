package repositoryUser

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *UsersRepository) CheckExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var user models.User
	err := r.DB.WithContext(ctx).First(&user, &models.User{ID: id}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

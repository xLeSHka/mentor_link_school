package repositoryUser

import (
	"context"

	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

func (r *UsersRepository) GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error) {
	err = r.DB.WithContext(ctx).First(&person, &models.User{ID: id}).Error
	return
}

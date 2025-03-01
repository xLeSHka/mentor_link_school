package repositoryUser

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) Login(ctx context.Context, person *models.User) (*models.User, error) {

	err := r.DB.WithContext(ctx).FirstOrCreate(person).Error
	if err != nil {
		return nil, err
	}
	return person, err
}

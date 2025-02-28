package repositoryUser

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"gorm.io/gorm/clause"
)

func (r *UsersRepository) Create(ctx context.Context, person *models.User) (*models.User, error) {
	err := r.DB.Create(person).WithContext(ctx).Clauses(clause.Returning{}).Error
	return person, err
}

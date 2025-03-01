package repositoryUser

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"gorm.io/gorm/clause"
)

func (r *UsersRepository) Login(ctx context.Context, person *models.User) (*models.User, error) {
	tx := r.DB.Begin()
	err := tx.FirstOrCreate(person).WithContext(ctx).Clauses(clause.Returning{}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return person, err
}

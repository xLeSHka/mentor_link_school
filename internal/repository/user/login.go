package repositoryUser

import (
	"context"
	"errors"

	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *UsersRepository) Login(ctx context.Context, person *models.User) (*models.User, error) {
	//err := r.DB.WithContext(ctx).FirstOrCreate(person).Error
	err := r.DB.WithContext(ctx).First(&person, &models.User{Name: person.Name}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.DB.WithContext(ctx).Create(person).Error; err != nil {
				return nil, err
			}
			return person, nil
		}
		return nil, err
	}
	return person, nil
}

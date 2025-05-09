package repositoryUser

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) Login(ctx context.Context, telegram string) (*models.User, error) {
	//err := r.DB.WithContext(ctx).FirstOrCreate(person).Error
	var person models.User
	err := r.DB.WithContext(ctx).Where("telegram = ?", telegram).First(&person).Error
	if err != nil {
		return nil, err
	}
	return &person, nil
}

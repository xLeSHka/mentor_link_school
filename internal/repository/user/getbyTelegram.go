package repositoryUser

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) GetByTelegram(ctx context.Context, telegram string) (person *models.User, err error) {
	err = r.DB.WithContext(ctx).First(&person, &models.User{Telegram: telegram}).Error
	return
}

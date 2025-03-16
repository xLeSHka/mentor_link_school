package repositoryUser

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"github.com/google/uuid"
)

func (r *UsersRepository) GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error) {
	err = r.DB.WithContext(ctx).First(&person, &models.User{ID: id}).Error
	return
}

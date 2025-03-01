package repositoryUser

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	return r.DB.Create(request).Error
}

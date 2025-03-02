package repositoryUser

import (
	"context"

	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersRepository) AddRole(ctx context.Context, role *models.Role) error {
	err := r.DB.WithContext(ctx).Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

package repositoryUser

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) AddRole(ctx context.Context, role *models.Role) error {
	err := r.DB.WithContext(ctx).Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

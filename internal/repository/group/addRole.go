package group

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository/user"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *repositoryUser.UsersRepository) AddRole(ctx context.Context, role *models.Role) error {
	err := r.DB.WithContext(ctx).Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

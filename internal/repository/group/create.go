package group

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) Create(ctx context.Context, group *models.Group) error {
	err := r.DB.Create(group).WithContext(ctx).Error
	return err
}

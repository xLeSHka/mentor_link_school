package group

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) Create(ctx context.Context, group *models.Group) error {
	tx := r.DB.Begin()

	err := tx.Create(group).WithContext(ctx).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Create(&models.Role{Role: "owner", GroupID: group.ID, UserID: group.UserID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

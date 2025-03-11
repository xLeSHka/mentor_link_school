package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *GroupRepository) Create(ctx context.Context, group *models.Group, userID uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Create(group).WithContext(ctx).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Create(&models.Role{Role: "owner", GroupID: group.ID, UserID: userID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

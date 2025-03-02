package repositoryUser

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *UsersRepository) CheckGroupExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var group models.Group
	err := r.DB.WithContext(ctx).First(&group, &models.Group{ID: id}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UsersRepository) GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	var group models.Group
	err := r.DB.WithContext(ctx).First(&group, &models.Group{InviteCode: &inviteCode}).Error
	return &group, err
}

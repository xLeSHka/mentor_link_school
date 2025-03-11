package repositoryUser

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	var group models.Group
	err := r.DB.WithContext(ctx).First(&group, &models.Group{InviteCode: &inviteCode}).Error
	return &group, err
}

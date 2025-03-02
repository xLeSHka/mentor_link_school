package group

import (
	"context"

	"github.com/google/uuid"
)

func (r *GroupRepository) UpdateInviteCode(ctx context.Context, groupID uuid.UUID, inviteCode string) error {
	return r.DB.WithContext(ctx).Table("groups").Where("id = ?", groupID).Update("invite_code", inviteCode).Error
}

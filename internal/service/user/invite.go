package userService

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (s *UsersService) Invite(ctx context.Context, inviteCode string, userID uuid.UUID) (bool, error) {
	group, err := s.usersRepository.GetGroupByInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("invite code not found")
		}
		return false, err
	}
	err = s.usersRepository.AddRole(ctx, &models.Role{
		UserID:  userID,
		GroupID: group.ID,
		Role:    "student",
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return false, errors.New("user already in group")
		}
		return false, err
	}
	return true, nil
}

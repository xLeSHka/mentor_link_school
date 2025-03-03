package userService

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return s.usersRepository.GetGroupByInviteCode(ctx, inviteCode)
}
